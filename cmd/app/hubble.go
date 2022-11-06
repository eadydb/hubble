package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eadydb/hubble/cmd/app/cmd"
	"github.com/eadydb/hubble/pkg/output"
	"github.com/eadydb/hubble/pkg/output/log"
	shell "github.com/kballard/go-shellquote"
)

func Run(out, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	catchStackdumpRequests()
	catchCtrlC(cancel)

	c := cmd.NewHubbleCommand(out, stderr)

	if cmdLine := os.Getenv("HUBBLE_CMDLINE"); cmdLine != "" && len(os.Args) == 1 {
		parsed, err := shell.Split(cmdLine)
		if err != nil {
			return fmt.Errorf("HUBBLE_CMDLINE is invalid: %w", err)
		}
		// XXX logged before logrus.SetLevel is called in NewHubbleCommand's PersistentPreRunE
		log.Entry(ctx).Debugf("Retrieving command line from HUBBLE_CMDLINE: %q", parsed)
		c.SetArgs(parsed)
	}
	c, err := c.ExecuteContextC(ctx)
	if err != nil {
		err = extractInvalidUsageError(err)
		if errors.Is(err, context.Canceled) {
			log.Entry(ctx).Debugln("ignore error since context is cancelled:", err)
		} else if !cmd.ShouldSuppressErrorReporting(c) {
			// As we allow some color setup using CLI flags for the main run, we can't run SetupColors()
			// for the entire skaffold run here. It's possible SetupColors() was never called, so call it again
			// before we print an error to get the right coloring.
			errOut := output.SetupColors(context.Background(), stderr, output.DefaultColorCode, false)
			output.Red.Fprintln(errOut, err)
		}
	}
	return err
}
