package config

import (
	"context"
	"fmt"
	"io"
)

func Unset(ctx context.Context, out io.Writer, args []string) error {
	if err := unsetConfigValue(args[0]); err != nil {
		return err
	}

	logUnsetConfigForUser(out, args[0])
	return nil
}

func logUnsetConfigForUser(out io.Writer, key string) {
	if global {
		fmt.Fprintf(out, "unset global value %s\n", key)
	} else {
		fmt.Fprintf(out, "unset value %s for context %s\n", key, kubecontext)
	}
}

func unsetConfigValue(name string) error {
	return setConfigValue(name, "")
}
