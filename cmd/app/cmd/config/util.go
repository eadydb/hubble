package config

import (
	"context"
	"fmt"

	"github.com/eadydb/hubble/pkg/config"
	kctx "github.com/eadydb/hubble/pkg/kubernetes/context"
	"github.com/eadydb/hubble/pkg/output/log"
	"github.com/eadydb/hubble/pkg/util"
)

func resolveKubectlContext() {
	if kubecontext != "" {
		return
	}

	config, err := kctx.CurrentConfig()
	switch {
	case err != nil:
		log.Entry(context.TODO()).Warn("unable to retrieve current kubectl context, using global values")
		global = true
	case config.CurrentContext == "":
		log.Entry(context.TODO()).Info("no kubectl context currently set, using global values")
		global = true
	default:
		kubecontext = config.CurrentContext
	}
}

func getConfigForKubectxOrDefault() (*config.ContextConfig, error) {
	cfg, err := getConfigForKubectx()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = &config.ContextConfig{}
		if !global {
			cfg.Kubecontext = kubecontext
		}
	}

	return cfg, nil
}

func getConfigForKubectx() (*config.ContextConfig, error) {
	resolveKubectlContext()

	if kubecontext == "" && !global {
		return nil, fmt.Errorf("missing `--kube-context` or `--global`")
	}

	cfg, err := config.ReadConfigFile(configFile)
	if err != nil {
		return nil, err
	}
	if global {
		return cfg.Global, nil
	}
	for _, contextCfg := range cfg.ContextConfigs {
		if util.RegexEqual(contextCfg.Kubecontext, kubecontext) {
			return contextCfg, nil
		}
	}
	return nil, nil
}
