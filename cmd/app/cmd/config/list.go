package config

import (
	"context"
	"fmt"
	"io"

	"github.com/eadydb/hubble/pkg/config"
	"github.com/eadydb/hubble/pkg/yaml"
)

func List(ctx context.Context, out io.Writer) error {
	var configYaml []byte
	if showAll {
		cfg, err := config.ReadConfigFile(configFile)
		if err != nil {
			return err
		}
		if cfg == nil || (cfg.Global == nil && len(cfg.ContextConfigs) == 0) { // empty config
			return nil
		}
		configYaml, err = yaml.Marshal(&cfg)
		if err != nil {
			return fmt.Errorf("marshaling config: %w", err)
		}
	} else {
		contextConfig, err := getConfigForKubectx()
		if err != nil {
			return err
		}
		if contextConfig == nil { // empty config
			return nil
		}
		configYaml, err = yaml.Marshal(&contextConfig)
		if err != nil {
			return fmt.Errorf("marshaling config: %w", err)
		}
	}

	fmt.Fprintf(out, "skaffold config: %s\n", configFile)
	out.Write(configYaml)

	return nil
}
