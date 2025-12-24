package client // import "go.microcore.dev/framework/transport/http/client"

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/http/client/core"
)

type Option func(*client)

func WithCore(core *fasthttp.Client) Option {
	return func(s *client) {
		s.core = core
	}
}

func WithCoreOptions(opts ...core.Option) Option {
	return func(s *client) {
		s.core = core.New(opts...)
	}
}

func WithTelemetryManager(telemetry telemetry.Manager) Option {
	return func(s *client) {
		s.telemetry = telemetry
	}
}

type RequestOption func(*request)

func WithRequestContext(context context.Context) RequestOption {
	return func(r *request) {
		r.context = context
	}
}

func WithRequestMethod(method string) RequestOption {
	return func(r *request) {
		r.method = method
	}
}

func WithRequestBody(body []byte) RequestOption {
	return func(r *request) {
		r.body = body
	}
}

func WithRequestJsonBody(data any) RequestOption {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(
			"failed to parse json body",
			slog.Any("error", err),
		)
		panic(err)
	}
	return func(r *request) {
		r.body = bytes
	}
}

func WithRequestHeaders(headers ...requestHeader) RequestOption {
	return func(r *request) {
		r.headers = headers
	}
}
