package core // import "go.microcore.dev/framework/transport/http/server/core"

import (
	"time"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
)

type Option func(*fasthttp.Server)

// Server name for sending in response headers.
func WithName(name string) Option {
	return func(s *fasthttp.Server) {
		s.Name = name
	}
}

// The maximum number of concurrent connections the server may serve.
func WithConcurrency(concurrency int) Option {
	return func(s *fasthttp.Server) {
		s.Concurrency = concurrency
	}
}

// Per-connection buffer size for requests' reading.
// This also limits the maximum header size.
//
// Increase this buffer if your clients send multi-KB RequestURIs
// and/or multi-KB headers (for example, BIG cookies).
func WithReadBufferSize(readBufferSize int) Option {
	return func(s *fasthttp.Server) {
		s.ReadBufferSize = readBufferSize
	}
}

// Per-connection buffer size for responses writing.
func WithWriteBufferSize(writeBufferSize int) Option {
	return func(s *fasthttp.Server) {
		s.WriteBufferSize = writeBufferSize
	}
}

// ReadTimeout is the amount of time allowed to read
// the full request including body. The connection's read
// deadline is reset when the connection opens, or for
// keep-alive connections after the first byte has been read.
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(s *fasthttp.Server) {
		s.ReadTimeout = readTimeout
	}
}

// WriteTimeout is the maximum duration before timing out
// writes of the response. It is reset after the request handler
// has returned.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(s *fasthttp.Server) {
		s.WriteTimeout = writeTimeout
	}
}

// IdleTimeout is the maximum amount of time to wait for the
// next request when keep-alive is enabled. If IdleTimeout
// is zero, the value of ReadTimeout is used.
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(s *fasthttp.Server) {
		s.IdleTimeout = idleTimeout
	}
}

// Maximum number of concurrent client connections allowed per IP.
func WithMaxConnsPerIP(maxConnsPerIP int) Option {
	return func(s *fasthttp.Server) {
		s.MaxConnsPerIP = maxConnsPerIP
	}
}

// Maximum number of requests served per connection.
//
// The server closes connection after the last request.
// 'Connection: close' header is added to the last response.
func WithMaxRequestsPerConn(maxRequestsPerConn int) Option {
	return func(s *fasthttp.Server) {
		s.MaxRequestsPerConn = maxRequestsPerConn
	}
}

// Whether to disable keep-alive connections.
//
// The server will close all the incoming connections after sending
// the first response to client if this option is set to true.
func WithMaxRequestBodySize(maxRequestBodySize int) Option {
	return func(s *fasthttp.Server) {
		s.MaxRequestBodySize = maxRequestBodySize
	}
}

// Maximum request body size.
// The server rejects requests with bodies exceeding this limit.
func WithDisableKeepalive(disableKeepalive bool) Option {
	return func(s *fasthttp.Server) {
		s.DisableKeepalive = disableKeepalive
	}
}

// Whether to enable tcp keep-alive connections.
// Whether the operating system should send tcp keep-alive messages
// on the tcp connection.
func WithTCPKeepalive(tcpKeepalive bool) Option {
	return func(s *fasthttp.Server) {
		s.TCPKeepalive = tcpKeepalive
	}
}

// Logs all errors, including the most frequent
// 'connection reset by peer', 'broken pipe' and 'connection timeout'
// errors. Such errors are common in production serving real-world
// clients.
func WithLogAllErrors(logAllErrors bool) Option {
	return func(s *fasthttp.Server) {
		s.LogAllErrors = logAllErrors
	}
}
