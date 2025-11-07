package core // import "go.microcore.dev/framework/transport/http/server/core"

import (
	"time"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/transport/http/server/core"

	DefaultServerName               = "microcore"
	DefaultServerConcurrency        = 256 * 1024
	DefaultServerReadBufferSize     = 4096 // 4 KB
	DefaultServerWriteBufferSize    = 4096 // 4 KB
	DefaultServerReadTimeout        = 10 * time.Second
	DefaultServerWriteTimeout       = 10 * time.Second
	DefaultServerIdleTimeout        = 10 * time.Second
	DefaultServerMaxConnsPerIP      = 0               // unlimited
	DefaultServerMaxRequestsPerConn = 0               // unlimited
	DefaultServerMaxRequestBodySize = 4 * 1024 * 1024 // 4 MB
	DefaultServerDisableKeepalive   = false
	DefaultServerTCPKeepalive       = false
	DefaultServerLogAllErrors       = true
)
