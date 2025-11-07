package core // import "go.microcore.dev/framework/transport/http/client/core"

import (
	"log/slog"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

var logger = log.New(pkg)

func New(opts ...Option) *fasthttp.Client {
	core := &fasthttp.Client{
		Name:                          DefaultClientName,
		MaxConnsPerHost:               DefaultClientMaxConnsPerHost,
		MaxIdleConnDuration:           DefaultClientMaxIdleConnDuration,
		MaxConnDuration:               DefaultClientMaxConnDuration,
		MaxIdemponentCallAttempts:     DefaultClientMaxIdemponentCallAttempts,
		ReadBufferSize:                DefaultClientReadBufferSize,
		WriteBufferSize:               DefaultClientWriteBufferSize,
		ReadTimeout:                   DefaultClientReadTimeout,
		WriteTimeout:                  DefaultClientWriteTimeout,
		DisableHeaderNamesNormalizing: DefaultClientDisableHeaderNamesNormalizing,
		DisablePathNormalizing:        DefaultClientDisablePathNormalizing,
		Dial:                          defaultClientDial,
	}

	for _, opt := range opts {
		opt(core)
	}

	logger.Info(
		"created",
		slog.String("name", core.Name),
	)

	return core
}
