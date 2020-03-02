package pkg

import (
	"fmt"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	v1 "k8s.io/api/core/v1"
)

const (
	annotationGSMsecretID            = "jenkins-x.io/gsm-secret-id"
	annotationGSMKubernetesSecretKey = "jenkins-x.io/gsm-kubernetes-secret-key"
)

func (o watchOptions) populateSecret(secret *v1.Secret, projectID string) error {
	if secret.Annotations[annotationGSMsecretID] == "" {
		return nil
	}

	secretID := secret.Annotations[annotationGSMsecretID]

	secretValue, err := o.getGoogleSecretManagerSecret(secretID, projectID)
	if err != nil {
		return fmt.Errorf("failed to find secret id %s in Google Secrets Manager: %v", secretID, err)
	}

	// initialise if secret has no data
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

func (o watchOptions) getGoogleSecretManagerSecret(secretID, projectID string) ([]byte, error) {

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Retrieve the secret
	result, err := o.smClient.AccessSecretVersion(o.ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret with id %s", secretID)
	}

	return result.Payload.Data, nil
}
