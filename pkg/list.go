package pkg

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"

	"github.com/jenkins-x-labs/gsm-controller/pkg/shared"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func ListSecrets(projectID string) error {

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

	secretsList, err := opts.kubeclient.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to list secrets in namespace %s", namespace)
	}

	for _, secret := range secretsList.Items {
		err := opts.findSecretData(secret)
		if err != nil {
			return errors.Wrapf(err, "failed to populate secret %s", secret.Name)
		}
	}

	return nil
}
