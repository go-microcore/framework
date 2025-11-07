package router // import "go.microcore.dev/framework/transport/http/server/router"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/fasthttp/router"
)

var logger = log.New(pkg)

func New(opts ...Option) *router.Router {
	router := router.New()

	router.SaveMatchedRoutePath = DefaultRouterSaveMatchedRoutePath
	router.RedirectTrailingSlash = DefaultRouterRedirectTrailingSlash
	router.RedirectFixedPath = DefaultRouterRedirectFixedPath
	router.HandleMethodNotAllowed = DefaultRouterHandleMethodNotAllowed
	router.HandleOPTIONS = DefaultRouterHandleOPTIONS
	router.GlobalOPTIONS = defaultRouterGlobalOPTIONS
	router.NotFound = defaultRouterNotFound
	router.MethodNotAllowed = defaultRouterMethodNotAllowed
	router.PanicHandler = defaultRouterPanicHandler

	for _, opt := range opts {
		opt(router)
	}

	logger.Info(
		"router has been successfully created",
		slog.Bool("save_matched_route_path", router.SaveMatchedRoutePath),
		slog.Bool("redirect_trailing_slash", router.RedirectTrailingSlash),
		slog.Bool("redirect_fixed_path", router.RedirectFixedPath),
		slog.Bool("handle_method_not_allowed", router.HandleMethodNotAllowed),
	)

	return router
}
