package postgres // import "go.microcore.dev/framework/db/postgres"

import (
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/db/postgres/client"

	"gorm.io/gorm"
)

type Option func(*p)

func WithClient(client *gorm.DB) Option {
	return func(p *p) {
		p.client = client
	}
}

func WithClientOptions(opts ...client.Option) Option {
	return func(p *p) {
		p.client = client.New(opts...)
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(p *p) {
		p.shutdownTimeout = timeout
	}
}

func WithoutShutdownHandler() Option {
	return func(p *p) {
		p.shutdownHandler = false
	}
}
