package shared

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type factory struct {
	kubeConfigCache *string
}

// NewFactory creates a factory with the default Kubernetes resources defined
// if optionalClientConfig is nil, then flags will be bound to a new clientcmd.ClientConfig.
// if optionalClientConfig is not nil, then this factory will make use of it.
func NewFactory() Factory {
	f := &factory{}
	return f
}

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
