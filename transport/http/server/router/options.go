package router // import "go.microcore.dev/framework/transport/http/server/router"

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
)

type Option func(*router.Router)

// If enabled, adds the matched route path onto the ctx.UserValue context
// before invoking the handler.
// The matched route path is only added to handlers of routes that were
// registered when this option was enabled.
func WithSaveMatchedRoutePath(saveMatchedRoutePath bool) Option {
	return func(r *router.Router) {
		r.SaveMatchedRoutePath = saveMatchedRoutePath
	}
}

// Enables automatic redirection if the current route can't be matched but a
// handler for the path with (without) the trailing slash exists.
// For example if /foo/ is requested but a route only exists for /foo, the
// client is redirected to /foo with http status code 301 for GET requests
// and 308 for all other request methods.
func WithRedirectTrailingSlash(redirectTrailingSlash bool) Option {
	return func(r *router.Router) {
		r.RedirectTrailingSlash = redirectTrailingSlash
	}
}

// If enabled, the router tries to fix the current request path, if no
// handle is registered for it.
// First superfluous path elements like ../ or // are removed.
// Afterwards the router does a case-insensitive lookup of the cleaned path.
// If a handle can be found for this route, the router makes a redirection
// to the corrected path with status code 301 for GET requests and 308 for
// all other request methods.
// For example /FOO and /..//Foo could be redirected to /foo.
// RedirectTrailingSlash is independent of this option.
func WithRedirectFixedPath(redirectFixedPath bool) Option {
	return func(r *router.Router) {
		r.RedirectFixedPath = redirectFixedPath
	}
}

// If enabled, the router checks if another method is allowed for the
// current route, if the current request can not be routed.
// If this is the case, the request is answered with 'Method Not Allowed'
// and HTTP status code 405.
// If no other Method is allowed, the request is delegated to the NotFound
// handler.
func WithHandleMethodNotAllowed(handleMethodNotAllowed bool) Option {
	return func(r *router.Router) {
		r.HandleMethodNotAllowed = handleMethodNotAllowed
	}
}

// If enabled, the router automatically replies to OPTIONS requests.
// Custom OPTIONS handlers take priority over automatic replies.
func WithHandleOPTIONS(handleOPTIONS bool) Option {
	return func(r *router.Router) {
		r.HandleOPTIONS = handleOPTIONS
	}
}

// An optional fasthttp.RequestHandler that is called on automatic OPTIONS requests.
// The handler is only called if HandleOPTIONS is true and no OPTIONS
// handler for the specific path was set.
// The "Allowed" header is set before calling the handler.
func WithGlobalOPTIONS(globalOPTIONS fasthttp.RequestHandler) Option {
	return func(r *router.Router) {
		r.GlobalOPTIONS = globalOPTIONS
	}
}

// Configurable fasthttp.RequestHandler which is called when no matching route is
// found. If it is not set, default NotFound is used.
func WithNotFound(notFound fasthttp.RequestHandler) Option {
	return func(r *router.Router) {
		r.NotFound = notFound
	}
}

// Configurable fasthttp.RequestHandler which is called when a request
// cannot be routed and HandleMethodNotAllowed is true.
// If it is not set, ctx.Error with fasthttp.StatusMethodNotAllowed is used.
// The "Allow" header with allowed request methods is set before the handler
// is called.
func WithMethodNotAllowed(methodNotAllowed fasthttp.RequestHandler) Option {
	return func(r *router.Router) {
		r.MethodNotAllowed = methodNotAllowed
	}
}

// Function to handle panics recovered from http handlers.
// It should be used to generate a error page and return the http error code
// 500 (Internal Server Error).
// The handler can be used to keep your server from crashing because of
// unrecovered panics.
func WithPanicHandler(panicHandler func(*fasthttp.RequestCtx, interface{})) Option {
	return func(r *router.Router) {
		r.PanicHandler = panicHandler
	}
}
