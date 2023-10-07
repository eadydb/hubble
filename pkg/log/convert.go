package log

// The following is Level security definitions.
const (
	InfoLevelSecurity  = "info"
	DebugLevelSecurity = "debug"
	WarnLevelSecurity  = "warn"
	ErrorLevelSecurity = "error"
)

// ToKlogLevel maps the current logging level to a Klog level integer
func ToKlogLevel(level Level) int {
	if int(level) > 0 {
		return 0
	}
	return -int(level)
}

// ToLogSeverityLevel maps the current logging level to a severity level string
func ToLogSeverityLevel(level Level) string {
	switch {
	case level < LevelInfo:
		return DebugLevelSecurity
	case level < LevelWarn:
		return InfoLevelSecurity
	case level < LevelError:
		return WarnLevelSecurity
	default:
		return ErrorLevelSecurity
	}
}
