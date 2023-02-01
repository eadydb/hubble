package context

import (
	"context"
	"fmt"
	"github.com/eadydb/hubble/pkg/log"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sync"
)

// For testing
var (
	CurrentConfig = getCurrentConfig
)

var (
	kubeConfigOnce sync.Once
	kubeConfig     clientcmd.ClientConfig

	configureOnce  sync.Once
	kubeContext    string
	kubeConfigFile string
)

// ConfigureKubeConfig sets an override for the current context in the k8s config.
func ConfigureKubeConfig(cliKubeConfig, cliKubeContext string) {
	configureOnce.Do(func() {
		kubeContext = cliKubeContext
		kubeConfigFile = cliKubeConfig
		if kubeContext != "" {
			log.Entry(context.TODO()).Info("Activated kube-context %q", kubeContext)
		}
	})
}

// GetDefaultRestClientConfig returns a REST client config for API calls against the Kubernetes API.
// If ConfigureKubeConfig was called before, the CurrentContext will be overridden.
func GetDefaultRestClientConfig() (*restclient.Config, error) {
	return getRestClientConfig(kubeContext, kubeConfigFile)
}

// GetRestClientConfig returns a REST client config for API calls against the Kubernetes API for the given context.
func GetRestClientConfig(kubeContext string) (*restclient.Config, error) {
	return getRestClientConfig(kubeContext, kubeConfigFile)
}

// GetClusterInfo returns the Cluster information for the given kubeContext
func GetClusterInfo(kctx string) (*clientcmdapi.Cluster, error) {
	rawConfig, err := getCurrentConfig()
	if err != nil {
		return nil, err
	}
	c, found := rawConfig.Clusters[kctx]
	if !found {
		return nil, fmt.Errorf("failed to get cluster info for kubeContext: `%s`", kctx)
	}
	return c, nil
}

func getRestClientConfig(kctx string, kcfg string) (*restclient.Config, error) {
	logger := log.Entry(context.TODO())
	logger.Debug("getting client config for kubeContext: `%s`", kctx)

	rawConfig, err := getCurrentConfig()
	if err != nil {
		return nil, err
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(rawConfig, kctx, &clientcmd.ConfigOverrides{CurrentContext: kctx},
		clientcmd.NewDefaultClientConfigLoadingRules())
	restConfig, err := clientConfig.ClientConfig()
	if kctx == "" && kcfg == "" && clientcmd.IsEmptyConfig(err) {
		logger.Debug("no kube-context set and no kubeConfig found, attempting in-cluster config")
		restConfig, err := restclient.InClusterConfig()
		if err != nil {
			return restConfig, fmt.Errorf("error creating REST client config in-cluster: %w", err)
		}

		return restConfig, nil
	}
	if err != nil {
		return restConfig, fmt.Errorf("error creating REST client config for kubeContext %q: %w", kctx, err)
	}

	return restConfig, nil
}

// getCurrentConfig retrieves and caches the raw kubeConfig.
func getCurrentConfig() (clientcmdapi.Config, error) {
	kubeConfigOnce.Do(func() {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = kubeConfigFile
		kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{
			CurrentContext: kubeContext,
		})
	})
	cfg, err := kubeConfig.RawConfig()
	if kubeContext != "" {
		cfg.CurrentContext = kubeContext
	}
	return cfg, err
}
