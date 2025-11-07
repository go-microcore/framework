package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"os"
	"syscall"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/shutdown"
)

var (
	shutdownSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
)
