package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/db/postgres/client"

	defaultDsnHost            = "localhost"
	defaultDsnPort            = 5432
	defaultDsnUser            = "postgres"
	defaultDsnPassword        = ""
	defaultDsnDb              = "postgres"
	defaultDsnSsl             = "disable"
	defaultDsnSearchPath      = "public"
	defaultDsnApplicationName = "microcore"
)
