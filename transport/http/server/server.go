package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strconv"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/http/server/core"
	"go.microcore.dev/framework/transport/http/server/listener"
	"go.microcore.dev/framework/transport/http/server/router"

	"github.com/aquasecurity/table"
	fasthttpRouter "github.com/fasthttp/router"
	"github.com/liamg/tml"
	fastHttpSwagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
)

type Interface interface {
	SetListener(listener net.Listener) Interface
	SetCore(core *fasthttp.Server) Interface
	SetRouter(router *fasthttpRouter.Router) Interface
	UseTelemetry(telemetry telemetry.Interface) Interface
	EnableTLS(tls *TLS) Interface
	AddMiddleware(MiddlewareHandler) Interface
	AddRoute(opts ...RouteOption) Interface
	AddRouteGroup(opts ...RouteGroupOption) Interface
	UseCors(opts ...CorsOption) Interface
	UseSwagger() Interface
	UseProfiling() Interface
	Listen() <-chan error
	Up()
	GetShutdownTimeout() time.Duration
	GetShutdownHandler() bool
	Shutdown(ctx context.Context, reason string) error
	ShutdownHandler(sig os.Signal) error
}

type server struct {
	listener        net.Listener
	core            *fasthttp.Server
	router          *fasthttpRouter.Router
	middleware      middleware
	telemetry       telemetry.Interface
	tls             *TLS
	shutdownTimeout time.Duration
	shutdownHandler bool
}

type route struct {
	method      string
	path        string
	handler     func(*RequestContext)
	middlewares []MiddlewareHandler
}
type rawRoute struct {
	method  string
	path    string
	handler RequestHandler
}
type routeGroup struct {
	path        string
	middlewares []MiddlewareHandler
	rawRoutes   []rawRoute
	routeGroups []*routeGroup
}

type cors struct {
	origin  string
	methods string
	headers string
}

type TLS struct {
	Cert string
	Key  string
}

type middleware = []func(fasthttp.RequestHandler) fasthttp.RequestHandler

type RequestHandler func(*RequestContext)
type MiddlewareHandler func(RequestHandler) RequestHandler

var logger = log.New(pkg)

func New(opts ...Option) Interface {
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
		shutdown.AddHandler(server.ShutdownHandler)
		logger.Info("shutdown handler has been successfully registered")
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

func (s *server) SetListener(listener net.Listener) Interface {
	s.listener = listener
	return s
}

func (s *server) SetCore(core *fasthttp.Server) Interface {
	s.core = core
	return s
}

func (s *server) SetRouter(router *fasthttpRouter.Router) Interface {
	s.router = router
	return s
}

func (s *server) UseTelemetry(telemetry telemetry.Interface) Interface {
	s.telemetry = telemetry
	return s
}

func (s *server) EnableTLS(tls *TLS) Interface {
	s.tls = tls
	return s
}

func (s *server) AddRoute(opts ...RouteOption) Interface {
	applyRoute(s.router, newRawRoute(opts...))
	return s
}

func (s *server) AddMiddleware(middleware MiddlewareHandler) Interface {
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

func (s *server) AddRouteGroup(opts ...RouteGroupOption) Interface {
	applyRouteGroup(s.router, nil, newRouteGroup(opts...), []MiddlewareHandler{})
	return s
}

func (s *server) UseCors(opts ...CorsOption) Interface {
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
				ctx.Response.Header.Set("Access-Control-Allow-Origin", cors.origin)
				ctx.Response.Header.Set("Access-Control-Allow-Methods", cors.methods)
				ctx.Response.Header.Set("Access-Control-Allow-Headers", cors.headers)
				handler(ctx)
			}
		},
	)

	return s
}

func (s *server) UseSwagger() Interface {
	s.router.Handle("GET", "/swagger/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
	return s
}

func (s *server) UseProfiling() Interface {
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

		strategies := map[bool]struct {
			mode  string
			tls   string
			serve func() error
		}{
			true: {
				mode: "network",
				tls:  "<red>Disable</red>",
				serve: func() error {
					return s.core.Serve(s.listener)
				},
			},
			false: {
				mode: "network (TLS)",
				tls:  "<green>Enable</green>",
				serve: func() error {
					return s.core.ServeTLS(s.listener, s.tls.Cert, s.tls.Key)
				},
			},
		}

		strategy := strategies[s.tls == nil]

		table := table.New(os.Stdout)
		table.SetHeaders("HTTP server")
		table.SetHeaderColSpans(0, 2)
		table.AddRow("Network", network)
		table.AddRow("Host", host)
		table.AddRow("Port", strconv.Itoa(port))
		table.AddRow("TLS", tml.Sprintf(strategy.tls))
		if s.telemetry == nil {
			table.AddRow("Telemetry", tml.Sprintf("<red>Disable</red>"))
		} else {
			table.AddRow("Telemetry", tml.Sprintf("<green>Enable</green>"))
		}
		table.Render()

		logger.Info(
			"started",
			slog.String("network", network),
			slog.String("host", host),
			slog.Int("port", port),
			slog.Bool("tls", s.tls != nil),
			slog.Bool("telemetry", s.telemetry != nil),
		)

		if err := strategy.serve(); err != nil {
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
	logger.Info(
		"shutting down",
		slog.String("reason", reason),
	)
	return s.core.ShutdownWithContext(ctx)
}

func (s *server) ShutdownHandler(sig os.Signal) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	reason := "unknown"
	if sig != nil {
		reason = sig.String()
	}

	return s.Shutdown(ctx, reason)
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
	var ctx context.Context
	if c := c.UserValue("ctx"); c != nil {
		ctx = c.(context.Context)
	} else {
		ctx = context.Background()
	}
	return ctx
}
