package postgres // import "go.microcore.dev/framework/db/postgres"

import (
	"time"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/db/postgres"

	DefaultShutdownTimeout = 10 * time.Second
	DefaultShutdownHandler = true
)
