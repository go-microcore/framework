package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"time"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/errors"
	"go.microcore.dev/framework/transport/http"
)

const (
	pkg = "go.microcore.dev/framework/transport/http/server"

	DefaultRouteMethod = http.MethodGet
	DefaultRoutePath   = "/"

	DefaultCorsOrigin  = "*"
	DefaultCorsMethods = "*"
	DefaultCorsHeaders = "*"

	DefaultShutdownTimeout = 10 * time.Second
	DefaultShutdownHandler = true
)

var (
	defaultResponseError = errors.ErrServiceUnavailable
)

func defaultRouteHandler(c *RequestContext) {
	c.Response.SetStatusCode(fasthttp.StatusOK)
	c.Response.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusOK))
}
