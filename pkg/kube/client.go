package kube

import (
	"io/fs"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// KubeClusterConfig kubernetes cluster config
type KubeClusterConfig struct {
	KubeConfigPath string // Config dir path
	Runner         Deploy // script deploy runner category
}

// KubeClusterClient kubernetes cluster client
type KubeClusterClient struct {
	Config *rest.Config
	Client *kubernetes.Clientset
}

type Deploy int

const (
	ClusterPodRunner Deploy = 1 // kubernetes cluster deploy probe
	AppRunner        Deploy = 2 // single deploy probe
)

// NewKubeClusterConfig init kubernetes cluster config
func NewKubeClusterConfig(path string, appRunner int) *KubeClusterConfig {
	runner := AppRunner
	if appRunner == 1 {
		runner = ClusterPodRunner
	}
	return &KubeClusterConfig{
		KubeConfigPath: path,
		Runner:         runner,
	}
}

// NewKubeClusterClient init kubernetes cluster client
func (kube *KubeClusterConfig) NewKubeClusterClient() *KubeClusterClient {
	var config *rest.Config
	var err error
	// out cluster client config
	if kube.Runner == AppRunner {
		kubeConfig := getClusterConfig(kube.KubeConfigPath)
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// in cluster client config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	clientside, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &KubeClusterClient{
		Config: config,
		Client: clientside,
	}
}

// getClusterConfig get kubernetes config path
func getClusterConfig(dir string) string {
	if fs.ValidPath(dir) {
		return dir
	}
	return path.Join(homedir.HomeDir(), ".kube", "config")
}
