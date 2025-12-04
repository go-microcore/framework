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

func New(opts ...Option) net.Listener {
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
		logger.Error(
			fmt.Sprintf("failed to create listener on %s://%s:%s",
				settings.network,
				settings.hostname,
				settings.port,
			),
			slog.Any("error", err),
		)
		panic(err)
	}

	logger.Debug(
		"listener has been successfully created",
		slog.String("network", settings.network),
		slog.String("hostname", settings.hostname),
		slog.String("port", settings.port),
	)

	return ln
}
