package constants

const (
	// These are phases in probe
	DevLoop     = Phase("DevLoop")
	Init        = Phase("Init")
	Build       = Phase("Build")
	Test        = Phase("Test")
	Render      = Phase("Render")
	Deploy      = Phase("Deploy")
	Verify      = Phase("Verify")
	StatusCheck = Phase("StatusCheck")
	PortForward = Phase("PortForward")
	Sync        = Phase("Sync")
	DevInit     = Phase("DevInit")
	Cleanup     = Phase("Cleanup")
)

const (
	Windows = "windows"

	// SubtaskIDNone is the value used for Event API messages when there is no
	// corresponding subtask
	SubtaskIDNone = "-1"
)

type Phase string
