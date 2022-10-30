package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Create creates or truncates a file to be used to output logs.
func Create(path ...string) (*os.File, error) {
	logfile := filepath.Join(os.TempDir(), "hubble")
	for _, p := range path {
		logfile = filepath.Join(logfile, escape(p))
	}

	dir := filepath.Dir(logfile)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("unable to create temp directory %q: %w", dir, err)
	}

	return os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
}

var escapeRegexp = regexp.MustCompile(`[^a-zA-Z0-9-_.]`)

func escape(s string) string {
	return escapeRegexp.ReplaceAllString(s, "-")
}
