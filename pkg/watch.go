package pkg

import (
	"fmt"

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

type secretOptions struct {
	kubeclient kubernetes.Interface
	projectID  string
}

func Foo(projectID string) error {

	f := shared.NewFactory()
	config, err := f.CreateKubeConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes config")
	}
	opts := secretOptions{
		projectID: projectID,
	}

	opts.kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes clientset")
	}

	factory := informers.NewSharedInformerFactory(opts.kubeclient, 0)

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
	err := PopulateSecret(secret, opts.projectID)
	if err != nil {
		klog.Error(err)
	}
	_, err = opts.kubeclient.CoreV1().Secrets(secret.Namespace).Update(secret)
	if err != nil {
		klog.Error(err)
	}
}

func (opts secretOptions) onUpdate(oldObj interface{}, newObj interface{}) {
	// only get the secret data from google manager store if the update even was because our annotation was added
	newSecret := newObj.(*v1.Secret)
	oldSecret := oldObj.(*v1.Secret)
	if oldSecret.Annotations[annotationGSMsecretID] == "" && newSecret.Annotations[annotationGSMsecretID] != "" {
		fmt.Println(newSecret.Name)
		err := PopulateSecret(newSecret, opts.projectID)
		if err != nil {
			klog.Error(err)
		}
		_, err = opts.kubeclient.CoreV1().Secrets(newSecret.Namespace).Update(newSecret)
		if err != nil {
			klog.Error(err)
		}
	}
}
