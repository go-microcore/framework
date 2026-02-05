package client // import "go.microcore.dev/framework/transport/http/client"

import (
	"context"
	"encoding/json"
	"fmt"

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

type RequestOption func(*request) error

func WithRequestContext(context context.Context) RequestOption {
	return func(r *request) error {
		r.context = context
		return nil
	}
}

func WithRequestMethod(method string) RequestOption {
	return func(r *request) error {
		r.method = method
		return nil
	}
}

func WithRequestBody(body []byte) RequestOption {
	return func(r *request) error {
		r.body = body
		return nil
	}
}

// WithRequestJsonBody serializes the given data to JSON and returns a RequestOption
// that sets it as the request body.
//
// If serialization fails, the function logs the error and immediately terminates
// the program with exit code 65 (ExitDataError), indicating invalid input data.
func WithRequestJsonBody(data any) RequestOption {
	return func(r *request) error {
		bytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to parse json body: %w", err)
		}
		r.body = bytes
		return nil
	}
}

func WithRequestHeaders(headers ...requestHeader) RequestOption {
	return func(r *request) error {
		r.headers = headers
		return nil
	}
}
