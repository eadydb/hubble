package event

import "github.com/sirupsen/logrus"

// logHook is an implementation of logrus.Hook used to send SkaffoldLogEvents
type logHook struct{}

func NewLogHook() logrus.Hook {
	return logHook{}
}

// Levels returns all levels as we want to send events for all levels
func (h logHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

// Fire constructs a SkaffoldLogEvent and sends it to the event channel
func (h logHook) Fire(entry *logrus.Entry) error {
	return nil
}
