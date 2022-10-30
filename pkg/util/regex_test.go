package util

import (
	"testing"

	"github.com/eadydb/hubble/pkg/testutil"
)

func TestKubectxEqual(t *testing.T) {
	ctxRe := ".*-i.*-am.*-test.*"
	ctxRePos := "wohoo-i-am-test-or"
	ctxReNeg := "test-am-i"
	testutil.CheckDeepEqual(t, true, RegexEqual(ctxRe, ctxRePos))
	testutil.CheckDeepEqual(t, false, RegexEqual(ctxRe, ctxReNeg))

	ctxStr := "^s^"
	ctxStrPos := "^s^"
	ctxStrNeg := "test-am-i"
	testutil.CheckDeepEqual(t, true, RegexEqual(ctxStr, ctxStrPos))
	testutil.CheckDeepEqual(t, false, RegexEqual(ctxStr, ctxStrNeg))
}
