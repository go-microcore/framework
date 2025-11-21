package client // import "go.microcore.dev/framework/db/redis/client"

import (
	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/db/redis/client"

	defaultNetwork    = "tcp"
	defaultAddr       = "localhost:6379"
	defaultClientName = "microcore"
	defaultDb         = 0
)
