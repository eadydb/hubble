package output

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/eadydb/hubble/pkg/constants"
	"github.com/eadydb/hubble/pkg/output/log"
)

const timestampFormat = "2006-01-02 15:04:05"

type hubbleWriter struct {
	MainWriter  io.Writer
	EventWriter io.Writer
	subtask     string

	timestamps bool
}

func (s hubbleWriter) Write(p []byte) (int, error) {
	written := 0
	if s.timestamps {
		t, err := s.MainWriter.Write([]byte(time.Now().Format(timestampFormat) + " "))
		if err != nil {
			return t, err
		}

		written += t
	}

	n, err := s.MainWriter.Write(p)
	if err != nil {
		return n, err
	}
	if n != len(p) {
		return n, io.ErrShortWrite
	}

	written += n

	s.EventWriter.Write(p)

	return written, nil
}

func GetWriter(ctx context.Context, out io.Writer, defaultColor int, forceColors bool, timestamps bool) io.Writer {
	if _, isSW := out.(hubbleWriter); isSW {
		return out
	}

	return hubbleWriter{
		MainWriter: SetupColors(ctx, out, defaultColor, forceColors),
		timestamps: timestamps,
	}
}

func IsStdout(out io.Writer) bool {
	sw, isSW := out.(hubbleWriter)
	if isSW {
		out = sw.MainWriter
	}
	cw, isCW := out.(colorableWriter)
	if isCW {
		out = cw.Writer
	}
	return out == os.Stdout
}

// GetUnderlyingWriter returns the underlying writer if out is a colorableWriter
func GetUnderlyingWriter(out io.Writer) io.Writer {
	sw, isSW := out.(hubbleWriter)
	if isSW {
		out = sw.MainWriter
	}
	cw, isCW := out.(colorableWriter)
	if isCW {
		out = cw.Writer
	}
	return out
}

// WithEventContext will return a new skaffoldWriter with the given parameters to be used for the event writer.
// If the passed io.Writer is not a skaffoldWriter, then it is simply returned.
func WithEventContext(ctx context.Context, out io.Writer, phase constants.Phase, subtaskID string) (io.Writer, context.Context) {
	ctx = context.WithValue(ctx, log.ContextKey, log.EventContext{
		Task:    phase,
		Subtask: subtaskID,
	})

	if sw, isSW := out.(hubbleWriter); isSW {
		return hubbleWriter{
			MainWriter: sw.MainWriter,
			subtask:    subtaskID,
			timestamps: sw.timestamps,
		}, ctx
	}

	return out, ctx
}
