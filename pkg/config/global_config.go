package config

// GlobalConfig holds configuration
type GlobalConfig struct {
	Global *ContextConfig `yaml:"global,omitempty"`
}

// ContextConfig
type ContextConfig struct {
	KubeConfig     *KubeConfig     `yaml:"kube"`
	RegistryConfig *RegistryConfig `yaml:"registry"`
}

// KubeConfig // kubernetes configuration
type KubeConfig struct {
	ConfigDir     string `yaml:"config_dir,omitempty"`
	ConfigFile    string `yaml:"config_file,omitempty"`
	Kubecontext   string `yaml:"kube-context,omitempty"`
	LocalCluster  string `yaml:"local-cluster"`
	KubeName      string `yaml:"kube-name,omitempty"`
	KubeVersion   string `yaml:"kube-version"`
	KubeApiServer string `yaml:"kube-api-server"`
}

// RegistryConfig registry configuration
type RegistryConfig struct {
	Registry string `yaml:"registry,omitempty"`
	Address  string `yaml:"address,omitempty"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
