package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"context"
	"encoding/json"
	"net"
	"time"

	fasthttpRouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	

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

func WithTelemetryManager(telemetry telemetry.Manager) Option {
	return func(s *server) {
		s.SetTelemetryManager(telemetry)
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

func WithRouteHandler(handler func(context.Context, *RequestContext)) RouteOption {
	return func(r *route) {
		r.handler = func(c *RequestContext) {
			handler(extractRequestContext(c.RequestCtx), c)
		}
	}
}

func WithRouteBodyParserHandler[T any](handler func(context.Context, *RequestContext, *T)) RouteOption {
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

			handler(extractRequestContext(c.RequestCtx), c, &body)
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
