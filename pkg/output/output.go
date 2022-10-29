package output

import "io"

const timestampFormat = "2006-01-02 15:04:05"

type probeWriter struct {
	MainWriter  io.Writer
	EventWriter io.Writer
	subtask     string

	timestamps bool
}

func (s probeWriter) Write(p []byte) (int, error) {
	return 0, nil
}
