package shared

import (
	"k8s.io/client-go/rest"
)

// Factory is the interface defined for Kubernetes, Jenkins X, and Tekton REST APIs
//go:generate pegomock generate github.com/jenkins-x/jx/pkg/jxfactory Factory -o mocks/factory.go
type Factory interface {
	// CreateKubeConfig creates the kubernetes configuration
	CreateKubeConfig() (*rest.Config, error)
}
