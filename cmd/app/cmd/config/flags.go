package config

import (
	"github.com/spf13/pflag"
)

var (
	configFile, kubecontext, surveyID string
	showAll, global, survey           bool
)

func AddCommonFlags(f *pflag.FlagSet) {
	f.StringVarP(&configFile, "config", "c", "", "Path to Skaffold config")
	f.StringVarP(&kubecontext, "kube-context", "k", "", "Kubectl context to set values against")
}

func AddListFlags(f *pflag.FlagSet) {
	f.BoolVarP(&showAll, "all", "a", false, "Show values for all kubecontexts")
}

func AddSetUnsetFlags(f *pflag.FlagSet) {
	f.BoolVarP(&global, "global", "g", false, "Set value for global config")
	f.BoolVarP(&survey, "survey", "s", false, "Set value for skaffold survey config")
	f.StringVarP(&surveyID, "id", "i", "", "Set value for given survey config")

	f.MarkHidden("survey")
	f.MarkHidden("id")
}
