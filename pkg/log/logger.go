package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
	"golang.org/x/term"
)

var IsTerminal = term.IsTerminal

func FromContext(ctx context.Context) *Logger {
	return wrapSlog(slog.FromContext(ctx))
}

func NewContext(ctx context.Context, logger *Logger) context.Context {
	return slog.NewContext(ctx, logger.log)
}

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
		ReplaceAttr: func(a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.AnyKind {
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
