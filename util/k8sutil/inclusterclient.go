package k8sutil // import "skyboat.io/x/util/k8sutil"

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// InClusterClient gets a local-cluster kubernetes client.
func InClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
