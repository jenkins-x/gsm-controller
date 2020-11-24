package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/joho/godotenv"
)

type secretOptions struct {
	kubeclient    kubernetes.Interface
	projectID     string
	accessSecrets accessSecrets
}

const (
	annotationGSMsecretID            = "jenkins-x.io/gsm-secret-id"
	annotationGSMKubernetesSecretKey = "jenkins-x.io/gsm-kubernetes-secret-key"
	annotationGSMSecretType          = "jenkins-x.io/gsm-type"
)

const (
	secretJsonType   = "json"
	secretDotenvType = "dotenv"
)

// New creates a instatialized Getter that can get files locally or remotely.
// useRemoteFS tells us if the service is configured to use the remote file system.
// accessKey and accessSecret are authentication parts for the remote file system.
func New(projectID string) *secretOptions {
	return &secretOptions{
		projectID:     projectID,
		accessSecrets: &googleSecretsManagerWrapper{},
	}
}

// minioWrapper adheres to the remoteFetcher interface
type googleSecretsManagerWrapper struct {
	smClient *secretmanager.Client
	ctx      context.Context
}

type accessSecrets interface {
	getGoogleSecretManagerSecret(secretID, projectID string) ([]byte, error)
}

func (o secretOptions) populateSecret(secret v1.Secret, projectID string) (v1.Secret, bool, error) {
	if secret.Annotations[annotationGSMsecretID] == "" {
		return secret, false, nil
	}

	secretID := secret.Annotations[annotationGSMsecretID]

	secretValue, err := o.accessSecrets.getGoogleSecretManagerSecret(secretID, projectID)
	if err != nil {
		return secret, false, fmt.Errorf("failed to find secret id %s in Google Secrets Manager: %v", secretID, err)
	}

	// initialise if secret has no data
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	// Treat as JSON value and save all keys into k8s secret
	if secret.Annotations[annotationGSMSecretType] == secretJsonType {
		var secretMap map[string]interface{}
		err := json.Unmarshal([]byte(secretValue), &secretMap)
		if err != nil {
			return secret, false, fmt.Errorf("failed to decode JSON secret id %s in Google Secrets Manager", secretID)
		}

		for key, value := range secretMap {
			secret.Data[key] = []byte(value.(string))
		}

	} else if secret.Annotations[annotationGSMSecretType] == secretDotenvType {
		r := bytes.NewReader(secretValue)
		secretMap, err := godotenv.Parse(r)

		if err != nil {
			return secret, false, fmt.Errorf("failed to decode dotenv %s in Google Secrets Manager", secretID)
		}
		for key, value := range secretMap {
			secret.Data[key] = []byte(value)
		}

	} else if secret.Annotations[annotationGSMKubernetesSecretKey] != "" {
		secret.Data[secret.Annotations[annotationGSMKubernetesSecretKey]] = secretValue
	} else {
		// default to the gsm secret id
		secret.Data[secret.Annotations[annotationGSMsecretID]] = secretValue
	}
	return secret, true, nil
}

func (o googleSecretsManagerWrapper) getGoogleSecretManagerSecret(secretID, projectID string) ([]byte, error) {

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Retrieve the secret
	result, err := o.smClient.AccessSecretVersion(o.ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret with name %s, err: %v", name, err)
	}

	return result.Payload.Data, nil
}
