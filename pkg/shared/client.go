package shared

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// PodNamespaceFile the file path and name for pod namespace
	PodNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

type factory struct {
	kubeConfigCache *string
}

// NewFactory creates a factory with the default Kubernetes resources defined
func NewFactory() Factory {
	f := &factory{}
	return f
}

// CreateKubeConfig figures out the kubernetes config from environment variables or default locations whether in or out
// of cluster
func (f *factory) CreateKubeConfig() (*rest.Config, error) {
	masterURL := ""
	kubeConfigEnv := os.Getenv("KUBECONFIG")
	if kubeConfigEnv != "" {
		pathList := filepath.SplitList(kubeConfigEnv)
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{Precedence: pathList},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterURL}}).ClientConfig()
	}
	kubeconfig := f.createKubeConfigText()
	var config *rest.Config
	var err error
	if kubeconfig != nil {
		exists, err := fileExists(*kubeconfig)
		if err == nil && exists {
			// use the current context in kubeconfig
			config, err = clientcmd.BuildConfigFromFlags(masterURL, *kubeconfig)
			if err != nil {
				return nil, err
			}
		}
	}
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (f *factory) createKubeConfigText() *string {
	var kubeconfig *string
	if f.kubeConfigCache != nil {
		return f.kubeConfigCache
	}
	text := ""
	if home := homeDir(); home != "" {
		text = filepath.Join(home, ".kube", "config")
	}
	kubeconfig = &text
	f.kubeConfigCache = kubeconfig
	return kubeconfig
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, errors.Wrapf(err, "failed to check if file exists %s", path)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return h
}

// CurrentContext returns the current context
func CurrentContext(config *api.Config) *api.Context {
	if config != nil {
		name := config.CurrentContext
		if name != "" && config.Contexts != nil {
			return config.Contexts[name]
		}
	}
	return nil
}

// CurrentNamespace returns the current namespace in the context
func CurrentNamespace() string {
	config, _, err := LoadConfig()
	ctx := CurrentContext(config)
	if ctx != nil {
		n := ctx.Namespace
		if n != "" {
			return n
		}
	}
	// if we are in a pod lets try load the pod namespace file
	data, err := ioutil.ReadFile(PodNamespaceFile)
	if err == nil {
		n := string(data)
		if n != "" {
			return n
		}
	}
	return "default"
}

// LoadConfig loads the Kubernetes configuration
func LoadConfig() (*api.Config, *clientcmd.PathOptions, error) {
	po := clientcmd.NewDefaultPathOptions()
	if po == nil {
		return nil, po, fmt.Errorf("Could not find any default path options for the kubeconfig file usually found at ~/.kube/config")
	}
	config, err := po.GetStartingConfig()
	if err != nil {
		return nil, po, fmt.Errorf("Could not load the kube config file %s due to %s", po.GetDefaultFilename(), err)
	}
	return config, po, err
}
