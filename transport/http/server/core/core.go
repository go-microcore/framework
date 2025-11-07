package core // import "go.microcore.dev/framework/transport/http/server/core"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/valyala/fasthttp"
)

var logger = log.New(pkg)

func New(opts ...Option) *fasthttp.Server {
	core := &fasthttp.Server{
		Name:                  DefaultServerName,
		Concurrency:           DefaultServerConcurrency,
		ReadBufferSize:        DefaultServerReadBufferSize,
		WriteBufferSize:       DefaultServerWriteBufferSize,
		ReadTimeout:           DefaultServerReadTimeout,
		WriteTimeout:          DefaultServerWriteTimeout,
		IdleTimeout:           DefaultServerIdleTimeout,
		MaxConnsPerIP:         DefaultServerMaxConnsPerIP,
		MaxRequestsPerConn:    DefaultServerMaxRequestsPerConn,
		MaxRequestBodySize:    DefaultServerMaxRequestBodySize,
		DisableKeepalive:      DefaultServerDisableKeepalive,
		TCPKeepalive:          DefaultServerTCPKeepalive,
		LogAllErrors:          DefaultServerLogAllErrors,
		CloseOnShutdown:       true,
		NoDefaultServerHeader: true,
		NoDefaultContentType:  true,
		NoDefaultDate:         true,
	}

	for _, opt := range opts {
		opt(core)
	}

	logger.Info(
		"core has been successfully created",
		slog.String("name", core.Name),
		slog.Int("concurrency", core.Concurrency),
		slog.Int("read_buffer_size", core.ReadBufferSize),
		slog.Int("write_buffer_size", core.WriteBufferSize),
	)

	return core
}
