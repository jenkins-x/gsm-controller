package pkg

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	v1 "k8s.io/api/core/v1"
)

const (
	annotationGSMsecretID            = "jenkins-x.io/gsm-secret-id"
	annotationGSMKubernetesSecretKey = "jenkins-x.io/gsm-kubernetes-secret-key"
)

func PopulateSecret(secret *v1.Secret, projectID string) error {
	if secret.Annotations[annotationGSMsecretID] == "" {
		return nil
	}

	secretID := secret.Annotations[annotationGSMsecretID]

	secretValue, err := getGoogleSecretManagerSecret(secretID, projectID)
	if err != nil {
		return fmt.Errorf("failed to find secret id %s in Google Secrets Manager: %v", secretID, err)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	if secret.Annotations[annotationGSMKubernetesSecretKey] != "" {
		secret.Data[secret.Annotations[annotationGSMKubernetesSecretKey]] = secretValue
	} else {
		// default to the gsm secret id

		secret.Data[secret.Annotations[annotationGSMsecretID]] = secretValue
	}
	return nil
}

func getGoogleSecretManagerSecret(secretID, projectID string) ([]byte, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup client")
	}

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret with id %s", secretID)
	}

	return result.Payload.Data, nil
}
