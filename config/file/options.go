package file // import "go.microcore.dev/framework/config/file"

import (
	_ "go.microcore.dev/framework"
)

type Option func(*c)

func WithPath(path string) Option {
	return func(c *c) {
		c.path = path
	}
}

func WithFormat(format format) Option {
	return func(c *c) {
		c.format = format
	}
}

func WithOut(out any) Option {
	return func(c *c) {
		c.out = out
	}
}
