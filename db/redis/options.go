package redis // import "go.microcore.dev/framework/db/redis"

import (
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/db/redis/client"

	"github.com/redis/go-redis/v9"
)

type Option func(*r)

func WithClient(client *redis.Client) Option {
	return func(r *r) {
		r.client = client
	}
}

func WithClientOptions(opts ...client.Option) Option {
	return func(r *r) {
		r.client = client.New(opts...)
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(r *r) {
		r.shutdownTimeout = timeout
	}
}

func WithoutShutdownHandler() Option {
	return func(r *r) {
		r.shutdownHandler = false
	}
}
