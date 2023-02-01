package log

import (
	"context"
	"golang.org/x/exp/slog"
)

var noop = wrapSlog(slog.New(noopHandler{}))

type noopHandler struct{}

var _ slog.Handler = noopHandler{}

func (noopHandler) Enabled(context.Context, Level) bool {
	return false
}

func (noopHandler) Handle(r slog.Record) error {
	return nil
}

func (h noopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h noopHandler) WithGroup(name string) slog.Handler {
	return h
}
