package log

import (
	"context"
	"os"

	"github.com/spf13/pflag"
)

// InitFlags initializes the flags for the log.
func InitFlags(ctx context.Context, flags *pflag.FlagSet) (context.Context, *Logger) {
	var level levelFlagValue
	flags.VarP(&level, "v", "v", "number for the log level verbosity (DEBUG, INFO, WARN, ERROR) or (-4, 0, 4, 8)")
	_ = flags.Parse(os.Args[1:])
	l := Level(level)
	logger := NewLogger(os.Stderr, l)
	return NewContext(ctx, logger), logger
}

type levelFlagValue Level

func (l levelFlagValue) String() string {
	return Level(l).String()
}

func (l *levelFlagValue) Set(v string) error {
	n, err := ParseLevel(v)
	if err != nil {
		return err
	}
	*l = levelFlagValue(n)
	return nil
}

func (l levelFlagValue) Type() string {
	return "log-level"
}
