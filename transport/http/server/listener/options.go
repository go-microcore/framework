package listener // import "go.microcore.dev/framework/transport/http/server/listener"

import (
	_ "go.microcore.dev/framework"
)

type Option func(*settings)

func WithNetwork(network string) Option {
	return func(s *settings) {
		s.network = network
	}
}

func WithHostname(hostname string) Option {
	return func(s *settings) {
		s.hostname = hostname
	}
}

func WithPort(port string) Option {
	return func(s *settings) {
		s.port = port
	}
}
