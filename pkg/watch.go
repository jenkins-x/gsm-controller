package pkg

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"

	"k8s.io/klog"

	"github.com/jenkins-x-labs/gsm-controller/pkg/shared"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type watchOptions struct {
	kubeclient    kubernetes.Interface
	projectID     string
	accessSecrets accessSecrets
}

func WatchSecrets(projectID string) error {

	var err error
	opts := New(projectID)

	// Create the google secrets manager client.
	gsm := googleSecretsManagerWrapper{
		ctx: context.Background(),
	}

	gsm.smClient, err = secretmanager.NewClient(gsm.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to setup secrets manager client")
	}

	opts.accessSecrets = gsm

	f := shared.NewFactory()
	config, err := f.CreateKubeConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes config")
	}

	opts.kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes clientset")
	}

	namespace := shared.CurrentNamespace()
	factory := informers.NewSharedInformerFactoryWithOptions(opts.kubeclient, 0, informers.WithNamespace(namespace))

	informer := factory.Core().V1().Secrets().Informer()

	stopper := make(chan struct{})
	defer close(stopper)

	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    opts.onAdd,
		UpdateFunc: opts.onUpdate,
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		return errors.New("Timed out waiting for caches to sync")
	}
	<-stopper
	return nil
}

func (opts watchOptions) onAdd(obj interface{}) {
	// Cast the obj as node
	secret := obj.(*v1.Secret)
	derefrencedSecret := *secret
	opts.findSecretData(derefrencedSecret)
}

func (opts watchOptions) onUpdate(oldObj interface{}, newObj interface{}) {
	// only get the secret data from google manager store if the update even was because our annotation was added
	newSecret := newObj.(*v1.Secret)
	oldSecret := oldObj.(*v1.Secret)
	if oldSecret.Annotations[annotationGSMsecretID] == "" && newSecret.Annotations[annotationGSMsecretID] != "" {
		derefrencedSecret := *newSecret
		opts.findSecretData(derefrencedSecret)
	}
}

func (opts watchOptions) findSecretData(secret v1.Secret) {
	secret, update, err := opts.populateSecret(secret, opts.projectID)
	if err != nil {
		klog.Error(err)
	}
	if update {
		_, err = opts.kubeclient.CoreV1().Secrets(secret.Namespace).Update(&secret)
		if err != nil {
			klog.Error(err)
		}
		klog.Infof("updated secret %s", secret.Name)
	}
}
