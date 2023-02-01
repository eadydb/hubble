package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
	"golang.org/x/term"
)

// IsTerminal returns true if the given file descriptor is a terminal.
var IsTerminal = term.IsTerminal

// Entry returns the Logger associated with ctx, or the default logger.
func Entry(ctx context.Context) *Logger {
	return wrapSlog(fromContext(ctx))
}

// NewContext returns a new context with the given logger.
func NewContext(ctx context.Context, logger *Logger) context.Context {
	return newContext(ctx, logger.log)
}

// NewLogger returns a new Logger that writes to w.
func NewLogger(w io.Writer, level slog.Level) *Logger {
	if w == nil {
		return noop
	}

	if file, ok := w.(*os.File); ok {
		fd := int(file.Fd())
		if IsTerminal(fd) {
			return wrapSlog(slog.New(newCtlHandler(w, fd, level)))
		}
	}

	handler := slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindAny {
				if t, ok := a.Value.Any().(fmt.Stringer); ok {
					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(t.String()),
					}
				}
			}
			return a
		},
	}
	return wrapSlog(slog.New(handler.NewJSONHandler(w)))
}
