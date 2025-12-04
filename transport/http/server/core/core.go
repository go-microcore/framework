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

	logger.Debug(
		"core has been successfully created",
		slog.String("name", core.Name),
		slog.Int("concurrency", core.Concurrency),
		slog.Int("read_buffer_size", core.ReadBufferSize),
		slog.Int("write_buffer_size", core.WriteBufferSize),
		slog.Duration("read_timeout", core.ReadTimeout),
		slog.Duration("write_timeout", core.WriteTimeout),
		slog.Duration("idle_timeout", core.IdleTimeout),
		slog.Int("max_conns_per_ip", core.MaxConnsPerIP),
		slog.Int("max_requests_per_conn", core.MaxRequestsPerConn),
		slog.Int("max_request_body_size", core.MaxRequestBodySize),
		slog.Bool("disable_keepalive", core.DisableKeepalive),
		slog.Bool("tcp_keepalive", core.TCPKeepalive),
		slog.Bool("log_all_errors", core.LogAllErrors),
	)

	return core
}
