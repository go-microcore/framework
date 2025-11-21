package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	fasthttpRouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/errors"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/http/server/core"
	"go.microcore.dev/framework/transport/http/server/listener"
	"go.microcore.dev/framework/transport/http/server/router"
)

type Option func(*server)

func WithListener(listener net.Listener) Option {
	return func(s *server) {
		s.listener = listener
	}
}

func WithListenerOptions(opts ...listener.Option) Option {
	return func(s *server) {
		s.listener = listener.New(opts...)
	}
}

func WithCore(core *fasthttp.Server) Option {
	return func(s *server) {
		s.core = core
	}
}

func WithCoreOptions(opts ...core.Option) Option {
	return func(s *server) {
		s.core = core.New(opts...)
	}
}

func WithRouter(router *fasthttpRouter.Router) Option {
	return func(s *server) {
		s.router = router
	}
}

func WithRouterOptions(opts ...router.Option) Option {
	return func(s *server) {
		s.router = router.New(opts...)
	}
}

func WithTelemetryManager(telemetry telemetry.Interface) Option {
	return func(s *server) {
		s.telemetry = telemetry

		s.middleware = append(
			s.middleware,
			func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
				return func(c *fasthttp.RequestCtx) {
					start := time.Now()

					ctx := s.telemetry.GetPropagator().Extract(extractRequestContext(c), fasthttpRequestCtxHeaderCarrier{c})
					ctx, span := telemetry.GetTracer().Start(ctx, "incoming http request")
					defer span.End()
					c.SetUserValue("ctx", ctx)

					defer func() {
						if rec := recover(); rec != nil {
							span.RecordError(fmt.Errorf("panic: %v", rec))
							span.SetStatus(codes.Error, fmt.Sprintf("panic: %v", rec))
							panic(rec)
						}

						duration := time.Since(start).Seconds()
						statusCode := c.Response.StatusCode()

						span.SetAttributes(
							attribute.String("path", string(c.Path())),
							attribute.String("method", string(c.Method())),
							attribute.Int("status", statusCode),
							attribute.Float64("duration", duration),
						)

						switch {
						case statusCode >= 500:
							span.SetStatus(codes.Error, fmt.Sprintf("Server error: HTTP %d", statusCode))
						case statusCode >= 400:
							span.SetStatus(codes.Error, fmt.Sprintf("Client error: HTTP %d", statusCode))
						default:
							span.SetStatus(codes.Ok, fmt.Sprintf("HTTP %d", statusCode))
						}
					}()

					handler(c)
				}
			},
		)
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}

func WithoutShutdownHandler() Option {
	return func(s *server) {
		s.shutdownHandler = false
	}
}

func WithTLS(tls *TLS) Option {
	return func(s *server) {
		s.tls = tls
	}
}

type RouteOption func(*route)

func WithRouteMethod(method string) RouteOption {
	return func(r *route) {
		r.method = method
	}
}

func WithRoutePath(path string) RouteOption {
	return func(r *route) {
		r.path = path
	}
}

func WithRouteHandler(handler func(*RequestContext, context.Context)) RouteOption {
	return func(r *route) {
		r.handler = func(c *RequestContext) {
			handler(c, extractRequestContext(c.RequestCtx))
		}
	}
}

func WithRouteBodyParserHandler[T any](handler func(*RequestContext, context.Context, *T)) RouteOption {
	return func(r *route) {
		r.handler = func(c *RequestContext) {
			var body T
			if err := json.Unmarshal(c.Request.Body(), &body); err == nil {
				if v, ok := any(&body).(interface{ Validate() error }); ok {
					if err := v.Validate(); err != nil {
						c.WriteError(err)
						return
					}
				}
			} else {
				c.WriteError(errors.ErrUnsupportedMediaType)
				return
			}

			handler(c, extractRequestContext(c.RequestCtx), &body)
		}
	}
}

func WithRouteMiddlewares(middlewares ...MiddlewareHandler) RouteOption {
	return func(r *route) {
		r.middlewares = middlewares
	}
}

type RouteGroupOption func(*routeGroup)

func WithRouteGroupPath(path string) RouteGroupOption {
	return func(r *routeGroup) {
		r.path = path
	}
}

func WithRouteGroupMiddlewares(middlewares ...MiddlewareHandler) RouteGroupOption {
	return func(r *routeGroup) {
		r.middlewares = middlewares
	}
}

func WithRouteGroupRoute(opts ...RouteOption) RouteGroupOption {
	return func(r *routeGroup) {
		r.rawRoutes = append(r.rawRoutes, *newRawRoute(opts...))
	}
}

func WithRouteGroup(opts ...RouteGroupOption) RouteGroupOption {
	return func(r *routeGroup) {
		r.routeGroups = append(r.routeGroups, newRouteGroup(opts...))
	}
}

type CorsOption func(*cors)

func WithCorsOrigin(origin string) CorsOption {
	return func(c *cors) {
		c.origin = origin
	}
}

func WithCorsMethods(methods string) CorsOption {
	return func(c *cors) {
		c.methods = methods
	}
}

func WithCorsHeaders(headers string) CorsOption {
	return func(c *cors) {
		c.headers = headers
	}
}
