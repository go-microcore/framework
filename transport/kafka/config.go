package kafka // import "go.microcore.dev/framework/transport/kafka"

import (
	"time"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/transport/kafka"

	DefaultShutdownTimeout = 10 * time.Second
	DefaultShutdownHandler = true
)
