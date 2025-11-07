package client // import "go.microcore.dev/framework/transport/http/client"

import (
	"context"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/transport/http"
)

const (
	pkg = "go.microcore.dev/framework/transport/http/client"

	DefaultRequestMethod = http.MethodGet
)

var (
	defaultRequestContext = context.Background()
)
