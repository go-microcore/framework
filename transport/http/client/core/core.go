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

	logger.Debug(
		"core created",
		slog.String("name", core.Name),
		slog.Int("max_conns_per_host", core.MaxConnsPerHost),
		slog.Duration("max_idle_conn_duration", core.MaxIdleConnDuration),
		slog.Duration("max_conn_duration", core.MaxConnDuration),
		slog.Int("max_idemponent_call_attempts", core.MaxIdemponentCallAttempts),
		slog.Int("read_buffer_size", core.ReadBufferSize),
		slog.Int("write_buffer_size", core.WriteBufferSize),
		slog.Duration("read_timeout", core.ReadTimeout),
		slog.Duration("write_timeout", core.WriteTimeout),
		slog.Int("max_response_body_size", core.MaxResponseBodySize),
		slog.Duration("max_conn_wait_timeout", core.MaxConnWaitTimeout),
		slog.Bool("no_default_user_agent_header", core.NoDefaultUserAgentHeader),
		slog.Bool("dial_dual_stack", core.DialDualStack),
		slog.Bool("disable_header_names_normalizing", core.DisableHeaderNamesNormalizing),
		slog.Bool("disable_path_normalizing", core.DisablePathNormalizing),
		slog.Bool("stream_response_body", core.StreamResponseBody),
	)

	return core
}
