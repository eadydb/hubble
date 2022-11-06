package main

import (
	"context"
	"errors"
	"os"

	"github.com/eadydb/hubble/cmd/app"
)

func main() {
	var code int
	if err := app.Run(os.Stdout, os.Stderr); err != nil && !errors.Is(err, context.Canceled) {
		// ignore cancelled errors
		code = app.ExitCode(err)
	}
	os.Exit(code)
}
