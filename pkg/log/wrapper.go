package log

import (
	"golang.org/x/exp/slog"
)

// Level is the logging level.
type Level = slog.Level

// The following is Level definitions copied from slog.
const (
	DebugLevel Level = slog.LevelDebug
	InfoLevel  Level = slog.LevelInfo
	WarnLevel  Level = slog.LevelWarn
	ErrorLevel Level = slog.LevelError
)

func wrapSlog(log *slog.Logger) *Logger {
	return &Logger{log}
}

// Logger is a wrapper around slog.Logger.
type Logger struct {
	log *slog.Logger
}

// LogDepth logs a message with the given level and depth.
func (l *Logger) LogDepth(callDepth int, level Level, msg string, args ...any) {
	l.log.LogDepth(callDepth+1, level, msg, args...)
}

// Log logs a message with the given level.
func (l *Logger) Log(level Level, msg string, args ...any) {
	l.LogDepth(0, level, msg, args...)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...any) {
	l.LogDepth(0, DebugLevel, msg, args...)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, args ...any) {
	l.LogDepth(0, InfoLevel, msg, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...any) {
	l.LogDepth(0, WarnLevel, msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args[:len(args):len(args)], slog.Any(slog.ErrorKey, err))
	}
	l.log.LogDepth(0, ErrorLevel, msg, args...)
}

// With returns a new Logger that includes the given arguments.
func (l *Logger) With(args ...any) *Logger {
	return wrapSlog(l.log.With(args...))
}

// WithGroup returns a new Logger that starts a group. The keys of all
// attributes added to the Logger will be qualified by the given name.
func (l *Logger) WithGroup(name string) *Logger {
	return wrapSlog(l.log.WithGroup(name))
}
