package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientConfig gets a k8s config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// MustGetClientConfig gets a k8s config or panics on error
func MustGetClientConfig(kubeconfig string) *rest.Config {
	cfg, err := GetClientConfig(kubeconfig)
	if err != nil {
		panic(err)
	}
	return cfg
}

// MustGetKubeClientFromConfig creates a Kubernetes client from config or panics
func MustGetKubeClientFromConfig(config *rest.Config) kubernetes.Interface {
	return kubernetes.NewForConfigOrDie(config)
}
