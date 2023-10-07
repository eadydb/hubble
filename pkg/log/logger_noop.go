package log

import (
	"context"

	"golang.org/x/exp/slog" //nolint:depguard
)

var noop = wrapSlog(noopHandler{}, LevelInfo)

type noopHandler struct{}

var _ slog.Handler = noopHandler{}

func (noopHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (noopHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h noopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h noopHandler) WithGroup(name string) slog.Handler {
	return h
}
