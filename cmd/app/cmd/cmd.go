package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eadydb/hubble/pkg/config"
	"github.com/eadydb/hubble/pkg/event"
	"github.com/eadydb/hubble/pkg/output"
	"github.com/eadydb/hubble/pkg/output/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	kubectx "github.com/eadydb/hubble/pkg/kubernetes/context"
)

var (
	opts         config.HubbleOptions
	v            string
	defaultColor int
	forceColors  bool
	overwrite    bool
	interactive  bool
	timestamps   bool
)

// Annotation for commands that should allow post execution housekeeping messages like updates and surveys
const (
	HouseKeepingMessagesAllowedAnnotation = "skaffold_annotation_housekeeping_allowed"
)

func NewHubbleCommand(out, errOut io.Writer) *cobra.Command {
	updateMsg := make(chan string, 1)

	rootCmd := &cobra.Command{
		Use:           "hubble",
		Short:         "Hubble is a tool to interact with the Hubble API",
		Long:          "Hubble is a tool to interact with the Hubble API",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().SilenceUsage = true

			opts.Command = cmd.Name()
			// Don't redirect output for Cobra internal `__complete` and `__completeNoDesc` commands.
			// These are used for command completion and send debug messages on stderr.
			if cmd.Name() != cobra.ShellCompRequestCmd && cmd.Name() != cobra.ShellCompNoDescRequestCmd {
				out := output.GetWriter(context.Background(), out, defaultColor, forceColors, timestamps)
				cmd.Root().SetOutput(out)

				// Setup logs
				if err := setUpLogs(errOut, v, timestamps); err != nil {
					return err
				}
			}

			// Setup kubeContext and kubeConfig
			kubectx.ConfigureKubeConfig(opts.KubeConfig, opts.KubeContext)
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			select {
			case msg := <-updateMsg:
				if err := config.UpdateMsgDisplayed(opts.GlobalConfig); err != nil {
					log.Entry(context.TODO()).Debugf("could not update the 'last-prompted' config for 'update-config' section due to %s", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "%s\n", msg)
			default:
			}
		},
	}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands:",
            Commands: []*cobra.Command{
                {
                    Use:   "init",
                    Short: "Initialize a new project",
                    Long:  "Initialize a new project",
				}
		},
		{
			Message: "Hubble Commands:",
            Commands: []*cobra.Command{
                {
                    Use:   "hubble",
                    Short: "Hubble is a tool to interact with the Hubble API",
				}
		},
	}

	groups.Add(rootCmd)

	// other commands
	rootCmd.AddCommand(NewCmdConfig())

	templates.ActsAsRootCommand(rootCmd, nil, groups...)

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", log.DefaultLogLevel.String(), fmt.Sprintf("Log level: one of %v", log.AllLevels))
	rootCmd.PersistentFlags().IntVar(&defaultColor, "color", int(output.DefaultColorCode), "Specify the default output color in ANSI escape codes")
	rootCmd.PersistentFlags().BoolVar(&forceColors, "force-colors", false, "Always print color codes (hidden)")
	rootCmd.PersistentFlags().BoolVar(&interactive, "interactive", true, "Allow user prompts for more information")
	rootCmd.PersistentFlags().BoolVar(&timestamps, "timestamps", false, "Print timestamps in logs")
	rootCmd.PersistentFlags().MarkHidden("force-colors")

	setFlagsFromEnvVariables(rootCmd)
	return rootCmd
}

func NewCmdOptions() *cobra.Command {
	cmd := &cobra.Command{
		Use: "options",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	templates.UseOptionsTemplates(cmd)

	return cmd
}

// Each flag can also be set with an env variable whose name starts with `HUBBLE_`.
func setFlagsFromEnvVariables(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		envVar := FlagToEnvVarName(f)
		if val, present := os.LookupEnv(envVar); present {
			rootCmd.PersistentFlags().Set(f.Name, val)
		}
	})
	for _, cmd := range rootCmd.Commands() {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			// special case for backward compatibility.
			if f.Name == "namespace" {
				if val, present := os.LookupEnv("HUBBLE_DEPLOY_NAMESPACE"); present {
					log.Entry(context.TODO()).Warn("Using HUBBLE_DEPLOY_NAMESPACE env variable is deprecated. Please use HUBBLE_NAMESPACE instead.")
					cmd.Flags().Set(f.Name, val)
				}
			}

			envVar := FlagToEnvVarName(f)
			if val, present := os.LookupEnv(envVar); present {
				cmd.Flags().Set(f.Name, val)
			}
		})
	}
}

func FlagToEnvVarName(f *pflag.Flag) string {
	return fmt.Sprintf("HUBBLE_%s", strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_"))
}

func setUpLogs(stdErr io.Writer, level string, timestamp bool) error {
	return log.SetupLogs(stdErr, level, timestamp, event.NewLogHook())
}

func isHouseKeepingMessagesAllowed(cmd *cobra.Command) bool {
	if cmd.Annotations == nil {
		return false
	}
	return cmd.Annotations[HouseKeepingMessagesAllowedAnnotation] == fmt.Sprintf("%t", true)
}

func allowHouseKeepingMessages(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[HouseKeepingMessagesAllowedAnnotation] = fmt.Sprintf("%t", true)
}

func isQuietMode() bool {
	switch {
	case !interactive:
		log.Entry(context.TODO()).Debug("Update check prompt, survey prompt and telemetry prompt disabled in non-interactive mode")
		return true
	case quietFlag:
		log.Entry(context.TODO()).Debug("Update check prompt, survey prompt and telemetry prompt disabled in quiet mode")
		return true
	case analyze:
		log.Entry(context.TODO()).Debug("Update check prompt, survey prompt and telemetry prompt disabled when running `init --analyze`")
		return true
	default:
		return false
	}
}
