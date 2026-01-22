package router // import "go.microcore.dev/framework/transport/http/server/router"

import (
	"log/slog"
	"runtime/debug"

	"github.com/valyala/fasthttp"
	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/transport/http/server/router"

	DefaultRouterSaveMatchedRoutePath   = false
	DefaultRouterRedirectTrailingSlash  = true
	DefaultRouterRedirectFixedPath      = true
	DefaultRouterHandleMethodNotAllowed = true
	DefaultRouterHandleOPTIONS          = true
)

func defaultRouterGlobalOPTIONS(ctx *fasthttp.RequestCtx) {
}

func defaultRouterNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Error(
		fasthttp.StatusMessage(fasthttp.StatusNotFound), 
		fasthttp.StatusNotFound,
	)
}

func defaultRouterMethodNotAllowed(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
}

func defaultRouterPanicHandler(ctx *fasthttp.RequestCtx, r any) {
	logger.Error(
		"panic in handler",
		slog.Any("error", r),
		slog.String("stack", string(debug.Stack())),
	)
	ctx.Error(
		fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
		fasthttp.StatusInternalServerError,
	)
}
