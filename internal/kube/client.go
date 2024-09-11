package kube

import (
	"fmt"

	"github.com/murtaza-u/rinc/internal/conf"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient creates a new Kubernetes API server client using the provided
// configuration.
func NewClient(c conf.KubernetesClient) (*kubernetes.Clientset, error) {
	if c.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("in cluster config: %w", err)
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("creating new kube client: %w", err)
		}
		return client, nil
	}
	if c.Kubeconfig == "" {
		return nil, fmt.Errorf("either `InCluster` or `Kubeconfig` must be set")
	}
	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("config from kubeconfig: %w", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating new kube client: %w", err)
	}
	return client, nil
}
