package config

// KubeConfig kubernetes cluster configuration
type KubeConfig struct {
	ClusterName      string `yaml:"cluster-name,omitempty"` // kubernetes cluster name
	ConfigPath       string `yaml:"config-path"`            // kubernetes config path
	ApiServerContext string `yaml:"api-server-context"`     // kubernetes api server context
	ClusterVersion   string `yaml:"cluster-version"`        // kubernetes cluster version
}

var KConfig KubeConfig
