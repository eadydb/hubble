package client

import (
	"fmt"
	"github.com/eadydb/hubble/pkg/kubernetes/context"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Initialize all known client auth plugins
)

// for tests
var (
	Client        = getClientset
	DynamicClient = getDynamicClient
	DefaultClient = getDefaultClientset
)

func getClientset(kubeContext string) (kubernetes.Interface, error) {
	config, err := context.GetRestClientConfig(kubeContext)
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func getDynamicClient(kubeContext string) (dynamic.Interface, error) {
	config, err := context.GetRestClientConfig(kubeContext)
	if err != nil {
		return nil, fmt.Errorf("getting client config for dynamic client: %w", err)
	}
	return dynamic.NewForConfig(config)
}

func getDefaultClientset() (kubernetes.Interface, error) {
	config, err := context.GetDefaultRestClientConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}
