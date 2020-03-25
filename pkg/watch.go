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

func (opts secretOptions) onAdd(obj interface{}) {
	// Cast the obj as node
	secret := obj.(*v1.Secret)
	derefrencedSecret := *secret
	err := opts.findSecretData(derefrencedSecret)
	if err != nil {
		klog.Error(err)
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
			klog.Error(err)
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
		klog.Infof("updated secret %s", secret.Name)
	}
	return nil
}
