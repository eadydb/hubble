package log

import (
	"context"
	"os"

	"github.com/spf13/pflag"
)

// InitFlags initializes the flags for the log.
func InitFlags(ctx context.Context, flags *pflag.FlagSet) (context.Context, *Logger) {
	v := flags.IntP("v", "v", 0, "number for the log level verbosity")
	_ = flags.Parse(os.Args[1:])
	logger := NewLogger(os.Stdout, Level(*v))
	return NewContext(ctx, logger), logger
}
