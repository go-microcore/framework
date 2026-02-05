package listener // import "go.microcore.dev/framework/transport/http/server/listener"

import (
	"fmt"
	"log/slog"
	"net"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

type settings struct {
	network, hostname, port string
}

var logger = log.New(pkg)

func New(opts ...Option) (net.Listener, error) {
	settings := &settings{
		network:  DefaultListenerNetwork,
		hostname: DefaultListenerHostname,
		port:     DefaultListenerPort,
	}

	for _, opt := range opts {
		opt(settings)
	}

	ln, err := net.Listen(settings.network, net.JoinHostPort(settings.hostname, settings.port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	logger.Debug(
		"listener created",
		slog.String("network", settings.network),
		slog.String("hostname", settings.hostname),
		slog.String("port", settings.port),
	)

	return ln, nil
}
