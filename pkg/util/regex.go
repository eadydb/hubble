package util

import (
	"context"
	re "regexp"
	"strings"

	"github.com/eadydb/hubble/pkg/output/log"
)

// RegexEqual matches the string 'actual' against a regex compiled from 'expected'
// If 'expected' is not a valid regex, string comparison is used as fallback
func RegexEqual(expected, actual string) bool {
	if strings.HasPrefix(expected, "!") {
		notExpected := expected[1:]

		return !regexMatch(notExpected, actual)
	}

	return regexMatch(expected, actual)
}

func regexMatch(expected, actual string) bool {
	if actual == expected {
		return true
	}

	matcher, err := re.Compile(expected)
	if err != nil {
		log.Entry(context.TODO()).Infof("context activation criteria '%s' is not a valid regexp, falling back to string", expected)
		return false
	}

	return matcher.MatchString(actual)
}
