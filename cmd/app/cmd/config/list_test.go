package config

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/eadydb/hubble/pkg/config"
	"github.com/eadydb/hubble/pkg/testutil"
	"github.com/eadydb/hubble/pkg/util"
	"github.com/eadydb/hubble/pkg/yaml"
)

func TestList(t *testing.T) {
	const dummyContext = "dummy-context"

	tests := []struct {
		cfg            *config.GlobalConfig
		name           string
		kubecontext    string
		global         bool
		showAll        bool
		expectedOutput string
	}{
		{
			name:        "print configs of a single kube-context",
			kubecontext: "this_is_a_context",
			cfg: &config.GlobalConfig{
				ContextConfigs: []*config.ContextConfig{
					{
						Kubecontext:        "another_context",
						DefaultRepo:        "other-value",
						MultiLevelRepo:     util.Ptr(false),
						LocalCluster:       util.Ptr(false),
						InsecureRegistries: []string{"good.io", "better.io"},
					},
					{
						Kubecontext:        "this_is_a_context",
						DefaultRepo:        "value",
						MultiLevelRepo:     util.Ptr(true),
						LocalCluster:       util.Ptr(true),
						InsecureRegistries: []string{"bad.io", "worse.io"},
					},
				},
			},
			expectedOutput: `kube-context: this_is_a_context
default-repo: value
multi-level-repo: true
local-cluster: true
insecure-registries:
- bad.io
- worse.io
`,
		},
		{
			name:   "print global configs",
			global: true,
			cfg: &config.GlobalConfig{
				Global: &config.ContextConfig{
					DefaultRepo:        "default-repo-value",
					MultiLevelRepo:     util.Ptr(true),
					LocalCluster:       util.Ptr(true),
					InsecureRegistries: []string{"mediocre.io"},
				},
				ContextConfigs: []*config.ContextConfig{
					{
						Kubecontext: "this_is_a_context",
						DefaultRepo: "value",
					},
				},
			},
			expectedOutput: `default-repo: default-repo-value
multi-level-repo: true
local-cluster: true
insecure-registries:
- mediocre.io
`,
		},
		{
			name:    "show all",
			showAll: true,
			cfg: &config.GlobalConfig{
				Global: &config.ContextConfig{
					DefaultRepo:        "default-repo-value",
					MultiLevelRepo:     util.Ptr(true),
					LocalCluster:       util.Ptr(true),
					InsecureRegistries: []string{"mediocre.io"},
				},
				ContextConfigs: []*config.ContextConfig{
					{
						Kubecontext: "this_is_a_context",
						DefaultRepo: "value",
					},
				},
			},
			expectedOutput: `
global:
  default-repo: default-repo-value
  multi-level-repo: true
  local-cluster: true
  insecure-registries:
  - mediocre.io
kubeContexts:
- kube-context: this_is_a_context
  default-repo: value
`,
		},
		{
			name:        "config has no values for kubecontext",
			kubecontext: "context-without-config",
			cfg: &config.GlobalConfig{
				Global: &config.ContextConfig{
					DefaultRepo:        "default-repo-value",
					MultiLevelRepo:     util.Ptr(true),
					LocalCluster:       util.Ptr(true),
					InsecureRegistries: []string{"mediocre.io"},
				},
			},
		},
		{
			name:        "config has no values for global",
			kubecontext: "context-without-config",
			cfg: &config.GlobalConfig{
				ContextConfigs: []*config.ContextConfig{
					{
						Kubecontext: "this_is_a_context",
						DefaultRepo: "value",
					},
				},
			},
		},
		{
			name:        "show all with empty config",
			showAll:     true,
			kubecontext: "context-without-config",
			cfg:         &config.GlobalConfig{},
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.name, func(t *testutil.T) {
			// create new config file
			content, _ := yaml.Marshal(*test.cfg)
			cfg := t.TempFile("config", content)

			t.Override(&config.ReadConfigFile, config.ReadConfigFileNoCache)
			t.Override(&configFile, cfg)
			t.Override(&global, test.global)
			t.Override(&showAll, test.showAll)
			if test.kubecontext != "" {
				t.Override(&kubecontext, test.kubecontext)
			} else {
				t.Override(&kubecontext, dummyContext)
			}

			buf := &bytes.Buffer{}
			// list values
			err := List(context.Background(), buf)
			t.CheckNoError(err)

			if test.expectedOutput != "" && !strings.HasSuffix(buf.String(), test.expectedOutput) {
				t.Errorf("expecting output to contain\n\n%s\nbut found\n\n%s\ninstead", test.expectedOutput, buf.String())
			}
			if test.expectedOutput == "" && buf.String() != "" {
				t.Errorf("expecting no output but found\n\n%s", buf.String())
			}
		})
	}
}
