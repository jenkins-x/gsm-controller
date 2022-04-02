package pkg

import (
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/magiconair/properties/assert"
)

func Test_extractSecretID(t *testing.T) {

	const examplePayload = `{
		"name": "projects/827585297303/secrets/test-secret/versions/1",
		"createTime": "2022-03-27T21:14:28.798342Z",
		"state": "ENABLED",
		"replicationStatus": {
			"automatic": {}
		},
		"etag": "\"15db39ae618386\""
	}`

	var deliveryAttempt int = 1

	msg := pubsub.Message{
		ID:          "4257778813345930",
		Data:        []byte(examplePayload),
		PublishTime: time.Unix(1648928005, 0),
		Attributes: map[string]string{
			"dataFormat": "JSON_API_V1",
			"eventType":  "SECRET_VERSION_ADD",
			"secretId":   "projects/827585297303/secrets/test-secret",
			"timestamp":  "2022-03-27T14:31:37.964139-07:00",
			"versionId":  "projects/827585297303/secrets/test-secret/versions/2",
		},
		DeliveryAttempt: &deliveryAttempt,
	}

	secretID, err := extractSecretIDIfNewVersion(&msg)
	assert.Equal(t, secretID, "test-secret")
	assert.Equal(t, err, nil, "err is nil")
}

func Test_extractSecretID_notAddingVersion(t *testing.T) {

	var deliveryAttempt int = 1
	msg := pubsub.Message{
		ID:          "4257778813345930",
		Data:        []byte(`{}`),
		PublishTime: time.Unix(1648928005, 0),
		Attributes: map[string]string{
			"dataFormat": "JSON_API_V1",
			"eventType":  "SECRET_DELETE",
			"secretId":   "projects/827585297303/secrets/test-secret",
			"timestamp":  "2022-03-27T14:31:37.964139-07:00",
		},
		DeliveryAttempt: &deliveryAttempt,
	}

	secretID, err := extractSecretIDIfNewVersion(&msg)
	assert.Equal(t, secretID, "")
	assert.Equal(t, err, nil, "err is nil")
}
