package redis // import "go.microcore.dev/framework/db/redis"

import (
	"time"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/db/redis"

	DefaultShutdownTimeout = 10 * time.Second
	DefaultShutdownHandler = true
)
