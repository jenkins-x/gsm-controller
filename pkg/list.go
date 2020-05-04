package pkg

import (
	"context"

	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"

	"github.com/jenkins-x-labs/gsm-controller/pkg/shared"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// ListOptions are the flags for list commands
type ListOptions struct {
	Cmd           *cobra.Command
	Args          []string
	projectID     string
	allNamespaces bool
}

var (
	list_desc = "List kubernetes secrets and update their data values from Google Secret Manager"

	list_example = "gsm list --project-id my-cool-gcp-project"
)

// NewCmdList creates the command
func NewCmdList() *cobra.Command {
	options := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   list_desc,
		Long:    list_desc,
		Example: list_example,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			shared.CheckErr(err)
		},
		SuggestFor: []string{"list"},
	}

	cmd.Flags().StringVarP(&options.projectID, "project-id", "", "", "The Google Project ID that contains the Google Secret Manager service")
	cmd.Flags().BoolVarP(&options.allNamespaces, "all-namespaces", "", false, "Scan all namespaces")

	return cmd
}

func (o ListOptions) Run() error {

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

	secretsList, err := secretOptions.kubeclient.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to list secrets in namespace %s", namespace)
	}

	for _, secret := range secretsList.Items {
		err := secretOptions.findSecretData(secret)
		if err != nil {
			return errors.Wrapf(err, "failed to populate secret %s", secret.Name)
		}
	}

	return nil
}
