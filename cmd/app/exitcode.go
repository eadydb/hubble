package app

import (
	"errors"
	"regexp"
)

type ExitCoder interface {
	ExitCode() int
}

// ExitCode extracts the exit code from the error.
func ExitCode(err error) int {
	var exitCoder ExitCoder
	if errors.As(err, &exitCoder) {
		return exitCoder.ExitCode()
	}
	return 1
}

type invalidUsageError struct{ err error }

func (i invalidUsageError) Unwrap() error { return i.err }
func (i invalidUsageError) Error() string { return i.err.Error() }
func (i invalidUsageError) ExitCode() int { return 127 }

// compiled list of common validation error prefixes from cobra/args.go and cobra/command.go based on skaffold's usage
var cobraUsageErrorPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^unknown command`),
	regexp.MustCompile(`^unknown( shorthand)? flag`),
	regexp.MustCompile(`^flag needs an argument:`),
	regexp.MustCompile(`^invalid argument `),
	regexp.MustCompile(`^accepts.*, received `),
}

func extractInvalidUsageError(err error) error {
	if err == nil {
		return nil
	}
	for _, pattern := range cobraUsageErrorPatterns {
		if pattern.MatchString(err.Error()) {
			return invalidUsageError{err}
		}
	}
	return err
}
