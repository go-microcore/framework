package client // import "go.microcore.dev/framework/transport/http/client"

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	netUrl "net/url"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/http/client/core"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Interface interface {
	SetCore(core *fasthttp.Client) Interface
	SetTelemetryManager(telemetry telemetry.Interface) Interface
	Request(url string, opts ...RequestOption) (*response, error)
}

type client struct {
	core      *fasthttp.Client
	telemetry telemetry.Interface
}

type request struct {
	context context.Context
	method  string
	body    []byte
	headers []requestHeader
}

type requestHeader struct {
	key   string
	value string
}

var logger = log.New(pkg)

func New(opts ...Option) Interface {
	client := &client{}

	for _, opt := range opts {
		opt(client)
	}

	if client.core == nil {
		client.core = core.New()
	}

	logger.Info(
		"created",
		slog.Bool("telemetry", client.telemetry != nil),
	)

	return client
}

func NewRequestHeader(key, value string) requestHeader {
	return requestHeader{key, value}
}

func (c *client) SetCore(core *fasthttp.Client) Interface {
	c.core = core
	return c
}

func (c *client) SetTelemetryManager(telemetry telemetry.Interface) Interface {
	c.telemetry = telemetry
	return c
}

func (c *client) Request(url string, opts ...RequestOption) (*response, error) {
	request := &request{
		context: defaultRequestContext,
		method:  DefaultRequestMethod,
	}

	for _, opt := range opts {
		opt(request)
	}

	if url == "" {
		return nil, errors.New("blank url")
	} else {
		u, err := netUrl.ParseRequestURI(url)
		if err != nil {
			return nil, fmt.Errorf("invalid url: %w", err)
		}
		if u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("url must include scheme and host")
		}
	}

	resCh := make(chan *response, 1)
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("request error: %v", r)
			}
		}()
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(url)
		req.Header.SetMethod(request.method)
		var span trace.Span
		if c.telemetry != nil {
			request.context, span = c.telemetry.GetTracer().Start(request.context, "outgoing http request")
			defer span.End()
			span.SetAttributes(
				attribute.String("request.method", request.method),
				attribute.String("request.url", url),
				attribute.Int("request.body_size", len(request.body)),
			)
			c.telemetry.GetPropagator().Inject(request.context, fasthttpRequestHeaderCarrier{&req.Header})
		}
		req.SetBodyRaw(request.body)
		for _, header := range request.headers {
			req.Header.Set(header.key, header.value)
		}
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)
		err := c.core.Do(req, resp)
		if err != nil {
			if c.telemetry != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			errCh <- err
			return
		}
		if c.telemetry != nil {
			span.SetAttributes(
				attribute.Int("response.status_code", resp.StatusCode()),
			)
		}
		cResp := &fasthttp.Response{}
		resp.CopyTo(cResp)
		resCh <- &response{cResp}
	}()
	select {
	case <-request.context.Done():
		return nil, fmt.Errorf("request failed due to context error: %w", request.context.Err())
	case r := <-resCh:
		return r, nil
	case e := <-errCh:
		return nil, e
	}
}

type fasthttpRequestHeaderCarrier struct {
	header *fasthttp.RequestHeader
}

func (f fasthttpRequestHeaderCarrier) Get(key string) string {
	return string(f.header.Peek(key))
}

func (f fasthttpRequestHeaderCarrier) Set(key, value string) {
	f.header.Set(key, value)
}

func (f fasthttpRequestHeaderCarrier) Keys() []string {
	keys := []string{}
	for k := range f.header.All() {
		keys = append(keys, string(k))
	}
	return keys
}
