package postgres // import "go.microcore.dev/framework/db/postgres"

import (
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/db/postgres/client"

	"gorm.io/gorm"
)

type Option func(*p) error

func WithClient(client *gorm.DB) Option {
	return func(p *p) error {
		p.client = client
		return nil
	}
}

func WithClientOptions(opts ...client.Option) Option {
	return func(p *p) error {
		client, err := client.New(opts...)
		if err != nil {
			return err
		}
		p.client = client
		return nil
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(p *p) error {
		p.shutdownTimeout = timeout
		return nil
	}
}

func WithoutShutdownHandler() Option {
	return func(p *p) error {
		p.shutdownHandler = false
		return nil
	}
}
