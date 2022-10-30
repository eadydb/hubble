package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	kubectx "github.com/eadydb/hubble/pkg/kubernetes/context"
	"github.com/eadydb/hubble/pkg/output/log"
	"github.com/eadydb/hubble/pkg/util"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigDir  = ".kube"
	defaultConfigFile = "config"
)

var (
	// config-file context
	configFile     *GlobalConfig
	configFileErr  error
	configFileOnce sync.Once

	// config context
	config     *ContextConfig
	configErr  error
	configOnce sync.Once

	ReadConfigFile             = readConfigFileCached
	GetConfigForCurrentKubectx = getConfigForCurrentKubectx

	current = time.Now()
)

// readConfigFileCached reads the specified file and returns the contents
// parsed into a GlobalConfig struct.
// This function will always return the identical data from the first read.
func readConfigFileCached(filename string) (*GlobalConfig, error) {
	configFileOnce.Do(func() {
		filenameOrDefault, err := ResolveConfigFile(filename)
		if err != nil {
			configFileErr = err
			log.Entry(context.TODO()).Warnf("Could not load global Skaffold defaults. Error resolving config file %q", filenameOrDefault)
			return
		}
		configFile, configFileErr = ReadConfigFileNoCache(filenameOrDefault)
		if configFileErr == nil {
			log.Entry(context.TODO()).Infof("Loaded Skaffold defaults from %q", filenameOrDefault)
		}
	})
	return configFile, configFileErr
}

// ReadConfigFileNoCache reads the given config yaml file and unmarshals the contents.
// Only visible for testing, use ReadConfigFile instead.
func ReadConfigFileNoCache(configFile string) (*GlobalConfig, error) {
	contents, err := os.ReadFile(configFile)
	if err != nil {
		log.Entry(context.TODO()).Warnf("Could not load global Skaffold defaults. Error encounter while reading file %q", configFile)
		return nil, fmt.Errorf("reading global config: %w", err)
	}
	config := GlobalConfig{}
	if err := yaml.Unmarshal(contents, &config); err != nil {
		log.Entry(context.TODO()).Warnf("Could not load global Skaffold defaults. Error encounter while unmarshalling the contents of file %q", configFile)
		return nil, fmt.Errorf("unmarshalling global skaffold config: %w", err)
	}
	return &config, nil
}

// ResolveConfigFile determines the default config location, if the configFile argument is empty.
func ResolveConfigFile(configFile string) (string, error) {
	if configFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("retrieving home directory: %w", err)
		}
		configFile = filepath.Join(home, defaultConfigDir, defaultConfigFile)
	}
	return configFile, util.VerifyOrCreateFile(configFile)
}

// GetConfigForCurrentKubectx returns the specific config to be modified based on the kubeContext.
// Either returns the config corresponding to the provided or current context,
// or the global config.
func getConfigForCurrentKubectx(configFile string) (*ContextConfig, error) {
	configOnce.Do(func() {
		cfg, err := ReadConfigFile(configFile)
		if err != nil {
			configErr = err
			return
		}
		kubeconfig, err := kubectx.CurrentConfig()
		if err != nil {
			configErr = err
			return
		}
		config, configErr = getConfigForKubeContextWithGlobalDefaults(cfg, kubeconfig.CurrentContext)
	})

	return config, configErr
}

func getConfigForKubeContextWithGlobalDefaults(cfg *GlobalConfig, kubeContext string) (*ContextConfig, error) {
	if kubeContext == "" {
		if cfg.Global == nil {
			return &ContextConfig{}, nil
		}
		return cfg.Global, nil
	}

	var mergedConfig ContextConfig
	for _, contextCfg := range cfg.ContextConfigs {
		if util.RegexEqual(contextCfg.KubeConfig.Kubecontext, kubeContext) {
			log.Entry(context.TODO()).Debugf("found config for context %q", kubeContext)
			mergedConfig = *contextCfg
		}
	}
	// in case there was no config for this kubeContext in cfg.ContextConfigs
	mergedConfig.KubeConfig.Kubecontext = kubeContext

	if cfg.Global != nil {
		// if values are unset for the current context, retrieve
		// the global config and use its values as a fallback.
		if err := mergo.Merge(&mergedConfig, cfg.Global, mergo.WithAppendSlice); err != nil {
			return nil, fmt.Errorf("merging context-specific and global config: %w", err)
		}
	}
	return &mergedConfig, nil
}
