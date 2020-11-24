package pkg

import (
	"errors"
	"testing"

	"github.com/magiconair/properties/assert"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeSecretsManager struct {
	data []byte
	err  error
}

func Test_watchOptions_populateSecret(t *testing.T) {

	secretWithAnnotation := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{annotationGSMsecretID: "foo"},
		},
	}

	successData := map[string][]byte{
		"foo": []byte("bar"),
	}
	type fields struct {
		data []byte
		err  error
	}
	tests := []struct {
		name         string
		fields       fields
		expectedData map[string][]byte
		wantErr      bool
	}{
		{name: "success", fields: struct {
			data []byte
			err  error
		}{data: []byte("bar"), err: nil}, expectedData: successData, wantErr: false},
		{name: "failed", fields: struct {
			data []byte
			err  error
		}{data: nil, err: errors.New("not found")}, expectedData: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := New("test_project")
			o.accessSecrets = &fakeSecretsManager{data: tt.fields.data, err: tt.fields.err}

			secretWithAnnotation, _, err := o.populateSecret(secretWithAnnotation, "test_project")

			if (err != nil) != tt.wantErr {
				t.Errorf("populateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, secretWithAnnotation.Data, tt.expectedData, "unexpected data set on secret")
		})
	}
}

func Test_watchOptions_populateJSONSecret(t *testing.T) {

	secretWithAnnotation := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{annotationGSMsecretID: "foo", annotationGSMSecretType: "json"},
		},
	}

	successData := map[string][]byte{
		"foo": []byte("bar"),
		"abc": []byte("def"),
	}
	type fields struct {
		data []byte
		err  error
	}
	tests := []struct {
		name         string
		fields       fields
		expectedData map[string][]byte
		wantErr      bool
	}{
		{name: "success", fields: struct {
			data []byte
			err  error
		}{data: []byte(`{"foo": "bar","abc": "def"}`), err: nil}, expectedData: successData, wantErr: false},

		{name: "bad_json", fields: struct {
			data []byte
			err  error
		}{data: []byte(`{"foo": bar`), err: nil}, expectedData: make(map[string][]byte), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := New("test_project")
			o.accessSecrets = &fakeSecretsManager{data: tt.fields.data, err: tt.fields.err}

			secretWithAnnotation, _, err := o.populateSecret(secretWithAnnotation, "test_project")

			if (err != nil) != tt.wantErr {
				t.Errorf("populateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, secretWithAnnotation.Data, tt.expectedData, "unexpected data set on secret")
		})
	}
}

func Test_watchOptions_populateDotenvSecret(t *testing.T) {

	secretWithAnnotation := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{annotationGSMsecretID: "foo", annotationGSMSecretType: "dotenv"},
		},
	}

	successData := map[string][]byte{
		"foo": []byte("bar"),
		"abc": []byte("def"),
	}
	type fields struct {
		data []byte
		err  error
	}
	tests := []struct {
		name         string
		fields       fields
		expectedData map[string][]byte
		wantErr      bool
	}{
		{name: "success", fields: struct {
			data []byte
			err  error
		}{data: []byte("foo=bar\nabc=def"), err: nil}, expectedData: successData, wantErr: false},

		{name: "bad_dotenv", fields: struct {
			data []byte
			err  error
		}{data: []byte("foo: \nbar"), err: nil}, expectedData: make(map[string][]byte), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := New("test_project")
			o.accessSecrets = &fakeSecretsManager{data: tt.fields.data, err: tt.fields.err}

			secretWithAnnotation, _, err := o.populateSecret(secretWithAnnotation, "test_project")

			if (err != nil) != tt.wantErr {
				t.Errorf("populateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, secretWithAnnotation.Data, tt.expectedData, "unexpected data set on secret")
		})
	}
}
func (f *fakeSecretsManager) getGoogleSecretManagerSecret(secretID, projectID string) ([]byte, error) {
	return f.data, f.err
}
