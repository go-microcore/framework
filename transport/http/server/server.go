package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"context"
	"log/slog"
	"net"
	"time"
	"fmt"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/http/server/core"
	"go.microcore.dev/framework/transport/http/server/listener"
	"go.microcore.dev/framework/transport/http/server/router"

	fasthttpRouter "github.com/fasthttp/router"
	fastHttpSwagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type (
	Manager interface {
		SetListener(listener net.Listener) Manager
		SetCore(core *fasthttp.Server) Manager
		SetRouter(router *fasthttpRouter.Router) Manager
		SetTelemetryManager(telemetry telemetry.Manager) Manager
		EnableTLS(tls *TLS) Manager
		AddMiddleware(MiddlewareHandler) Manager
		AddRoute(opts ...RouteOption) Manager
		AddRouteGroup(opts ...RouteGroupOption) Manager
		UseCors(opts ...CorsOption) Manager
		UseSwagger() Manager
		UseProfiling() Manager
		Listen() <-chan error
		Up()
		GetShutdownTimeout() time.Duration
		GetShutdownHandler() bool
		Shutdown(ctx context.Context, reason string) error
	}

	server struct {
		listener        net.Listener
		core            *fasthttp.Server
		router          *fasthttpRouter.Router
		middleware      middleware
		telemetry       telemetry.Manager
		tls             *TLS
		shutdownTimeout time.Duration
		shutdownHandler bool
	}

	route struct {
		method      string
		path        string
		handler     func(*RequestContext)
		middlewares []MiddlewareHandler
	}
	rawRoute struct {
		method  string
		path    string
		handler RequestHandler
	}
	routeGroup struct {
		path        string
		middlewares []MiddlewareHandler
		rawRoutes   []rawRoute
		routeGroups []*routeGroup
	}

	cors struct {
		origin  string
		methods string
		headers string
	}

	TLS struct {
		Cert string
		Key  string
	}

	middleware = []func(fasthttp.RequestHandler) fasthttp.RequestHandler

	RequestHandler    func(*RequestContext)
	MiddlewareHandler func(RequestHandler) RequestHandler
)

var logger = log.New(pkg)

func New(opts ...Option) Manager {
	server := &server{
		shutdownTimeout: DefaultShutdownTimeout,
		shutdownHandler: DefaultShutdownHandler,
	}

	for _, opt := range opts {
		opt(server)
	}

	if server.listener == nil {
		server.listener = listener.New()
	}

	if server.core == nil {
		server.core = core.New()
	}

	if server.router == nil {
		server.router = router.New()
	}

	if server.shutdownHandler {
		shutdown.AddHandler(server.Shutdown)
		logger.Debug("shutdown handler has been successfully registered")
	}

	logger.Info(
		"created",
		slog.Group("shutdown",
			slog.Duration("timeout", server.shutdownTimeout),
			slog.Bool("handler", server.shutdownHandler),
		),
		slog.Bool("telemetry", server.telemetry != nil),
		slog.Bool("tls", server.tls != nil),
	)

	return server
}

func (s *server) SetListener(listener net.Listener) Manager {
	s.listener = listener
	return s
}

func (s *server) SetCore(core *fasthttp.Server) Manager {
	s.core = core
	return s
}

func (s *server) SetRouter(router *fasthttpRouter.Router) Manager {
	s.router = router
	return s
}

func (s *server) SetTelemetryManager(telemetry telemetry.Manager) Manager {
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
	return s
}

func (s *server) EnableTLS(tls *TLS) Manager {
	s.tls = tls
	return s
}

func (s *server) AddRoute(opts ...RouteOption) Manager {
	applyRoute(s.router, newRawRoute(opts...))
	return s
}

func (s *server) AddMiddleware(middleware MiddlewareHandler) Manager {
	s.middleware = append(
		s.middleware,
		func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
			return func(ctx *fasthttp.RequestCtx) {
				middleware(
					func(ctx *RequestContext) {
						handler(ctx.RequestCtx)
					},
				)(&RequestContext{RequestCtx: ctx})
			}
		},
	)
	return s
}

func (s *server) AddRouteGroup(opts ...RouteGroupOption) Manager {
	applyRouteGroup(s.router, nil, newRouteGroup(opts...), []MiddlewareHandler{})
	return s
}

func (s *server) UseCors(opts ...CorsOption) Manager {
	cors := &cors{
		origin:  DefaultCorsOrigin,
		methods: DefaultCorsMethods,
		headers: DefaultCorsHeaders,
	}

	for _, opt := range opts {
		opt(cors)
	}

	s.middleware = append(
		s.middleware,
		func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
			return func(ctx *fasthttp.RequestCtx) {
				handler(ctx)
				ctx.Response.Header.Set("Access-Control-Allow-Origin", cors.origin)
				ctx.Response.Header.Set("Access-Control-Allow-Methods", cors.methods)
				ctx.Response.Header.Set("Access-Control-Allow-Headers", cors.headers)
			}
		},
	)

	return s
}

func (s *server) UseSwagger() Manager {
	s.router.Handle("GET", "/swagger/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
	return s
}

func (s *server) UseProfiling() Manager {
	s.router.Handle("GET", "/debug/pprof/{profile:*}", pprofhandler.PprofHandler)
	return s
}

func (s *server) Listen() <-chan error {
	exit := make(chan error, 1)
	go func() {
		s.middleware = append(
			s.middleware,
			func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
				return func(c *fasthttp.RequestCtx) {
					ctx := extractRequestContext(c)
					if s.telemetry != nil {
						ctx = s.telemetry.GetPropagator().Extract(ctx, fasthttpRequestCtxHeaderCarrier{c})
					}
					defer func() {
						logger.LogAttrs(
							ctx,
							slog.LevelInfo,
							"request",
							slog.String("method", string(c.Method())),
							slog.String("path", string(c.Path())),
							slog.Int("status", c.Response.StatusCode()),
						)
					}()
					handler(c)
				}
			},
		)

		handler := s.router.Handler
		for i := len(s.middleware) - 1; i >= 0; i-- {
			handler = s.middleware[i](handler)
		}
		s.core.Handler = handler

		addr := s.listener.Addr().(*net.TCPAddr)
		host := addr.IP.String()
		port := addr.Port
		network := addr.Network()

		serve := map[bool]func() error{
			true: func() error {
				return s.core.Serve(s.listener)
			},
			false: func() error {
				return s.core.ServeTLS(s.listener, s.tls.Cert, s.tls.Key)
			},
		}

		logger.Info(
			"started",
			slog.String("network", network),
			slog.String("host", host),
			slog.Int("port", port),
			slog.Bool("tls", s.tls != nil),
			slog.Bool("telemetry", s.telemetry != nil),
		)

		if err := serve[s.tls == nil](); err != nil {
			exit <- err
		}

		close(exit)
	}()
	return exit
}

func (s *server) Up() {
	if err := <-s.Listen(); err != nil {
		logger.Error(
			"failed to listen",
			slog.Any("error", err),
		)
	}
	logger.Info("stopped")
}

func (s *server) GetShutdownTimeout() time.Duration {
	return s.shutdownTimeout
}

func (s *server) GetShutdownHandler() bool {
	return s.shutdownHandler
}

func (s *server) Shutdown(ctx context.Context, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	logger.Debug(
		"shutdown",
		slog.String("reason", reason),
	)

	return s.core.ShutdownWithContext(ctx)
}

func newRawRoute(opts ...RouteOption) *rawRoute {
	route := &route{
		method:      DefaultRouteMethod,
		path:        DefaultRoutePath,
		handler:     defaultRouteHandler,
		middlewares: []MiddlewareHandler{},
	}

	for _, opt := range opts {
		opt(route)
	}

	handler := RequestHandler(route.handler)
	for i := len(route.middlewares) - 1; i >= 0; i-- {
		handler = route.middlewares[i](handler)
	}

	return &rawRoute{
		method:  route.method,
		path:    route.path,
		handler: handler,
	}
}

func applyRoute(router *fasthttpRouter.Router, route *rawRoute) {
	router.Handle(
		route.method,
		route.path,
		func(ctx *fasthttp.RequestCtx) {
			route.handler(
				&RequestContext{
					RequestCtx: ctx,
				},
			)
		},
	)
}

func newRouteGroup(opts ...RouteGroupOption) *routeGroup {
	routeGroup := &routeGroup{
		path:        DefaultRoutePath,
		middlewares: []MiddlewareHandler{},
		rawRoutes:   []rawRoute{},
		routeGroups: []*routeGroup{},
	}

	for _, opt := range opts {
		opt(routeGroup)
	}

	return routeGroup
}

func applyRouteGroup(router *fasthttpRouter.Router, group *fasthttpRouter.Group, routeGroup *routeGroup, middlewares []MiddlewareHandler) {
	if group == nil {
		group = router.Group(routeGroup.path)
	} else {
		group = group.Group(routeGroup.path)
	}

	routeGroup.middlewares = append(middlewares, routeGroup.middlewares...)

	for _, rawRoute := range routeGroup.rawRoutes {
		for i := len(routeGroup.middlewares) - 1; i >= 0; i-- {
			rawRoute.handler = routeGroup.middlewares[i](rawRoute.handler)
		}
		group.Handle(
			rawRoute.method,
			rawRoute.path,
			func(ctx *fasthttp.RequestCtx) {
				rawRoute.handler(
					&RequestContext{
						RequestCtx: ctx,
					},
				)
			},
		)
	}

	for _, g := range routeGroup.routeGroups {
		applyRouteGroup(router, group, g, routeGroup.middlewares)
	}
}

// This context is built through the entire middleware
// chain and can contain data for distributed tracing
// or any other implementations.
func extractRequestContext(c *fasthttp.RequestCtx) context.Context {
	if c := c.UserValue("ctx"); c != nil {
		return c.(context.Context)
	}
	return context.Background()
}
