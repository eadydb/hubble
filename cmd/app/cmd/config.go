package cmd

import (
	"github.com/eadydb/hubble/cmd/app/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Interact with the global Hubble config file (defaults to `$HOME/.hubble/config`)",
	}

	cmd.AddCommand(NewCmdSet())
	cmd.AddCommand(NewCmdUnset())
	cmd.AddCommand(NewCmdList())
	return cmd
}

func NewCmdSet() *cobra.Command {
	return NewCmd("set").
		WithDescription("Set a value in the global Hubble config").
		WithExample("Mark a registry as insecure", "config set insecure-registries <insecure1.io>").
		WithExample("Globally set the default image repository", "config set default-repo <myrepo>").
		WithExample("Globally set multi-level repo support", "config set multi-level-repo true").
		WithExample("Disable pushing images for a given Kubernetes context", "config set --kube-context <mycluster> local-cluster true").
		WithFlagAdder(func(f *pflag.FlagSet) {
			config.AddCommonFlags(f)
			config.AddSetUnsetFlags(f)
		}).
		ExactArgs(2, config.Set)
}

func NewCmdUnset() *cobra.Command {
	return NewCmd("unset").
		WithDescription("Unset a value in the global Hubble config").
		WithFlagAdder(func(f *pflag.FlagSet) {
			config.AddCommonFlags(f)
			config.AddSetUnsetFlags(f)
		}).
		ExactArgs(1, config.Unset)
}

func NewCmdList() *cobra.Command {
	return NewCmd("list").
		WithDescription("List all values set in the global Hubble config").
		WithFlagAdder(func(f *pflag.FlagSet) {
			config.AddCommonFlags(f)
			config.AddListFlags(f)
		}).
		NoArgs(config.List)
}
