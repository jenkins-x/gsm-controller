package pkg

import (
	"context"

	"github.com/spf13/cobra"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"

	"github.com/jenkins-x-labs/gsm-controller/pkg/shared"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// UpgradeOptions are the flags for delete commands
type WatchOptions struct {
	Cmd           *cobra.Command
	Args          []string
	projectID     string
	allNamespaces bool
}

var (
	watch_desc = "Watch kubernetes secrets and update their data values from Google Secret Manager"

	watch_example = "gsm watch --project-id my-cool-gcp-project"
)

// NewCmdWatch creates the command
func NewCmdWatch() *cobra.Command {
	options := &WatchOptions{}

	cmd := &cobra.Command{
		Use:     "watch",
		Long:    watch_desc,
		Short:   watch_desc,
		Example: watch_example,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			shared.CheckErr(err)
		},
		SuggestFor: []string{"watch"},
	}

	cmd.Flags().StringVarP(&options.projectID, "project-id", "", "", "The Google Project ID that contains the Google Secret Manager service")
	cmd.Flags().BoolVarP(&options.allNamespaces, "all-namespaces", "", false, "Scan all namespaces")
	return cmd
}

// Run implements this command
func (o *WatchOptions) Run() error {

	if o.projectID == "" {
		return errors.New("missing flag project-id")
	}

	var err error
	secretOptions := New(o.projectID)

	// Create the google secrets manager client.
	gsm := googleSecretsManagerWrapper{
		ctx: context.Background(),
	}

	gsm.smClient, err = secretmanager.NewClient(gsm.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to setup secrets manager client")
	}

	secretOptions.accessSecrets = gsm

	f := shared.NewFactory()
	config, err := f.CreateKubeConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes config")
	}

	secretOptions.kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes clientset")
	}

	namespace := shared.CurrentNamespace()
	if o.allNamespaces {
		namespace = v1.NamespaceAll
	}

	factory := informers.NewSharedInformerFactoryWithOptions(secretOptions.kubeclient, 0, informers.WithNamespace(namespace))

	informer := factory.Core().V1().Secrets().Informer()

	stopper := make(chan struct{})
	defer close(stopper)

	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    secretOptions.onAdd,
		UpdateFunc: secretOptions.onUpdate,
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		return errors.New("Timed out waiting for caches to sync")
	}
	<-stopper
	return nil
}

func (opts secretOptions) onAdd(obj interface{}) {
	// Cast the obj as node
	secret := obj.(*v1.Secret)
	derefrencedSecret := *secret
	err := opts.findSecretData(derefrencedSecret)
	if err != nil {
		log.Logger().Error(err)
	}
}

func (opts secretOptions) onUpdate(oldObj interface{}, newObj interface{}) {
	// only get the secret data from google manager store if the update even was because our annotation was added
	newSecret := newObj.(*v1.Secret)
	oldSecret := oldObj.(*v1.Secret)
	if oldSecret.Annotations[annotationGSMsecretID] == "" && newSecret.Annotations[annotationGSMsecretID] != "" {
		derefrencedSecret := *newSecret
		err := opts.findSecretData(derefrencedSecret)
		if err != nil {
			log.Logger().Error(err)
		}
	}
}

func (opts secretOptions) findSecretData(secret v1.Secret) error {
	secret, update, err := opts.populateSecret(secret, opts.projectID)
	if err != nil {
		return errors.Wrapf(err, "failed to populate secret %s from gsm in project %s", secret.Name, opts.projectID)
	}
	if update {
		_, err = opts.kubeclient.CoreV1().Secrets(secret.Namespace).Update(&secret)
		if err != nil {
			return errors.Wrapf(err, "failed to update secret %s", secret.Name)
		}
		log.Logger().Infof("updated secret %s", secret.Name)
	}
	return nil
}
