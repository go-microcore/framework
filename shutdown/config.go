package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"os"
	"syscall"
	"time"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/shutdown"

	defaultShutdownTimeout = 60 * time.Second
)

var (
	shutdownSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
)
