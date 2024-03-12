package pkg

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"cloud.google.com/go/pubsub"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"

	"github.com/jenkins-x-labs/gsm-controller/pkg/shared"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// SubscribeOptions are the flags for subscribe commands
type SubscribeOptions struct {
	Cmd           *cobra.Command
	Args          []string
	projectID     string
	subID         string
	allNamespaces bool
}

var (
	subscribe_desc = "Subscribe to a Cloud Pub/Sub topic and update kubernetes secrets when Google Secret Manager versions are added"

	subscribe_example = "gsm subscribe --project-id my-cool-gcp-project --subscription secret-events.gsm-controller"
)

// NewCmdSubscribe creates the comman
func NewCmdSubscribe() *cobra.Command {
	options := &SubscribeOptions{}

	cmd := &cobra.Command{
		Use:     "subscribe",
		Long:    subscribe_desc,
		Short:   subscribe_desc,
		Example: subscribe_example,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			shared.CheckErr(err)
		},
		SuggestFor: []string{"subscribe"},
	}

	cmd.Flags().StringVarP(&options.projectID, "project-id", "", "", "The Google Project ID that contains the Google Secret Manager service")
	cmd.Flags().StringVarP(&options.subID, "subscription", "", "", "The Google Pub/Sub subscription is set up to receive Google Secret Manager events")
	cmd.Flags().BoolVarP(&options.allNamespaces, "all-namespaces", "", false, "Scan all namespaces when looking for secret to update")
	return cmd
}

func (o SubscribeOptions) Run() error {

	if o.projectID == "" {
		return errors.New("missing flag project-id")
	}
	if o.subID == "" {
		return errors.New("missing flag subscription")
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

	secretOptions.namespace = shared.CurrentNamespace()
	if o.allNamespaces {
		secretOptions.namespace = v1.NamespaceAll
	}

	// Die early if listing will fail later
	_, err = secretOptions.kubeclient.CoreV1().Secrets(secretOptions.namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to list secrets in namespace %s", secretOptions.namespace)
	}

	ctx := context.Background() // FIXME
	pubsubClient, err := pubsub.NewClient(ctx, o.projectID)
	if err != nil {
		return errors.Wrap(err, "failed to create pubsub client")
	}
	defer pubsubClient.Close()

	sub := pubsubClient.Subscription(o.subID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		secretID, err := extractSecretIDIfNewVersion(msg)
		if err != nil {
			log.Logger().Error(err)
		}
		if secretID == "" {
			msg.Ack()
			return
		}

		err = secretOptions.handleVersionNotification(secretID)
		if err != nil {
			log.Logger().Error(err)
			if msg.DeliveryAttempt != nil && *msg.DeliveryAttempt < 3 {
				msg.Nack()
				return
			}
		}
		msg.Ack()
	})
	if err != nil {
		return errors.Wrap(err, "I want to come back to this error message")
	}
	return nil
}

func extractSecretIDIfNewVersion(msg *pubsub.Message) (string, error) {
	eventType, _ := msg.Attributes["eventType"]
	if eventType != "SECRET_VERSION_ADD" {
		return "", nil
	}
	fullSecretID, ok := msg.Attributes["secretId"]
	if !ok {
		log.Logger().Warning("Received SECRET_VERSION_ADD with no secretId")
		return "", nil
	}

	var projectNumber int
	var secretID string
	_, err := fmt.Sscanf(fullSecretID, "projects/%d/secrets/%s", &projectNumber, &secretID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to extract final portion of secretId %q\n", fullSecretID)
	}
	return secretID, nil
}

func (opts secretOptions) handleVersionNotification(gsmSecretID string) error {

	secret, err := opts.findSecret(gsmSecretID)
	if err != nil {
		return errors.Wrapf(err, "failed finding secret %q", gsmSecretID)
	}
	if secret == nil {
		log.Logger().Infof("No matching secret in kubernetes for %q\n", gsmSecretID)
		return nil
	}
	return opts.findSecretData(*secret)
}

func (opts secretOptions) findSecret(gsmSecretID string) (*v1.Secret, error) {

	secretsList, err := opts.kubeclient.CoreV1().Secrets(opts.namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list secrets in namespace %s", opts.namespace)
	}
	var _ = secretsList
	for _, secretItem := range secretsList.Items {
		anno, ok := secretItem.Annotations[annotationGSMsecretID]
		if !ok {
			continue
		}
		if anno == gsmSecretID {
			return &secretItem, nil
		}
	}

	return nil, nil
}
