package app

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func catchStackdumpRequests() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGUSR1,
	)

	go func() {
		for {
			<-signals
			buf := make([]byte, 1<<20)
			runtime.Stack(buf, true)
			os.Stderr.Write([]byte(fmt.Sprintf("---stacktraces begin: %v---\n", os.Args)))
			os.Stderr.Write(buf)
			os.Stderr.Write([]byte(fmt.Sprintf("---stacktraces end: %v---\n", os.Args)))
		}
	}()
}
