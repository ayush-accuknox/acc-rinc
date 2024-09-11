package conf

// KubernetesClient contains the configuration needed to communicate with the
// Kubernetes API server.
type KubernetesClient struct {
	// InCluster, when set to true, attempts to authenticate with the API
	// server using a service account token.
	//
	// Either `InCluster` must be set to true or the path to a kubeconfig file
	// must be provided below.
	InCluster bool `koanf:"inCluster"`
	// Kubeconfig is the path to the `kubeconfig` file. This is useful when
	// running the application outside the cluster.
	//
	// Either `InCluster` must be set to true or the path to a kubeconfig file
	// must be provided here.
	Kubeconfig string `koanf:"kubeconfig"`
}
