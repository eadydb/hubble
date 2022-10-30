package term

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/eadydb/hubble/pkg/constants"
	"github.com/eadydb/hubble/pkg/util"
	"golang.org/x/term"
)

func IsTerminal(w io.Writer) (uintptr, bool) {
	type descriptor interface {
		Fd() uintptr
	}

	if f, ok := w.(descriptor); ok {
		termFd := f.Fd()
		isTerm := term.IsTerminal(int(termFd))
		return termFd, isTerm
	}

	return 0, false
}

func SupportsColor(ctx context.Context) (bool, error) {
	if runtime.GOOS == constants.Windows {
		return true, nil
	}

	cmd := exec.Command("tput", "colors")
	res, err := util.RunCmdOut(ctx, cmd)
	if err != nil {
		return false, fmt.Errorf("checking terminal colors: %w", err)
	}

	numColors, err := strconv.Atoi(strings.TrimSpace(string(res)))
	if err != nil {
		return false, err
	}

	return numColors > 0, nil
}
