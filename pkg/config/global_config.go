package config

// GlobalConfig holds configuration
type GlobalConfig struct {
	Global         *ContextConfig   `yaml:"global,omitempty"`
	ContextConfigs []*ContextConfig `yaml:"kubeContexts"`
}

// ContextConfig is the context-specific config information provided in
// the global Skaffold config.
type ContextConfig struct {
	Kubecontext        string   `yaml:"kube-context,omitempty"`
	DefaultRepo        string   `yaml:"default-repo,omitempty"`
	MultiLevelRepo     *bool    `yaml:"multi-level-repo,omitempty"`
	LocalCluster       *bool    `yaml:"local-cluster,omitempty"`
	InsecureRegistries []string `yaml:"insecure-registries,omitempty"`
	// DebugHelpersRegistry is the registry from which the debug helper images are used.
	DebugHelpersRegistry string        `yaml:"debug-helpers-registry,omitempty"`
	UpdateCheck          *bool         `yaml:"update-check,omitempty"`
	Survey               *SurveyConfig `yaml:"survey,omitempty"`
	KindDisableLoad      *bool         `yaml:"kind-disable-load,omitempty"`
	K3dDisableLoad       *bool         `yaml:"k3d-disable-load,omitempty"`
	CollectMetrics       *bool         `yaml:"collect-metrics,omitempty"`
	UpdateCheckConfig    *UpdateConfig `yaml:"update,omitempty"`
	Registry             *Registry     `yaml:"registry,omitempty"`
}

// SurveyConfig is the survey config information
type SurveyConfig struct {
	DisablePrompt *bool         `yaml:"disable-prompt,omitempty"`
	LastTaken     string        `yaml:"last-taken,omitempty"`
	LastPrompted  string        `yaml:"last-prompted,omitempty"`
	UserSurveys   []*UserSurvey `yaml:"user-surveys,omitempty"`
}

type UserSurvey struct {
	ID    string `yaml:"id"`
	Taken *bool  `yaml:"taken,omitempty"`
}

// UpdateConfig is the update config information
type UpdateConfig struct {
	// TODO (tejaldesai) Move ContextConfig.UpdateCheck config within this struct
	LastPrompted string `yaml:"last-prompted,omitempty"`
}

// Registry is the server registry information
type Registry struct {
	Eureka *Eureka `yaml:"eureka,omitempty"`
	Nacos  *Nacos  `yaml:"nacos,omitempty"`
	Apollo *Apollo `yaml:"apollo,omitempty"`
}

type Eureka struct {
	Url      string `yaml:"url,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Nacos struct {
	Url       string `yaml:"url,omitempty"`
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

type Apollo struct {
	Url   string `yaml:"server,omitempty"`
	Token string `yaml:"token,omitempty"`
}
