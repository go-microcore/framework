package client // import "go.microcore.dev/framework/transport/http/client"

import (
	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
)

type response struct {
	*fasthttp.Response
}
