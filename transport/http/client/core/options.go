package core // import "go.microcore.dev/framework/transport/http/client/core"

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
)

type Option func(*fasthttp.Client)

// Transport defines a transport-like mechanism that wraps every request/response.
func WithTransport(transport fasthttp.RoundTripper) Option {
	return func(c *fasthttp.Client) {
		c.Transport = transport
	}
}

// Callback for establishing new connections to hosts.
//
// Default DialTimeout is used if not set.
func WithDialTimeout(dialTimeout fasthttp.DialFuncWithTimeout) Option {
	return func(c *fasthttp.Client) {
		c.DialTimeout = dialTimeout
	}
}

// Callback for establishing new connections to hosts.
//
// Note that if Dial is set instead of DialTimeout, Dial will ignore Request timeout.
// If you want the tcp dial process to account for request timeouts, use DialTimeout instead.
//
// If not set, DialTimeout is used.
func WithDial(dial fasthttp.DialFunc) Option {
	return func(c *fasthttp.Client) {
		c.Dial = dial
	}
}

// TLS config for https connections.
//
// Default TLS config is used if not set.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(c *fasthttp.Client) {
		c.TLSConfig = tlsConfig
	}
}

// When the client encounters an error during a request, the behavior—whether to retry
// and whether to reset the request timeout—should be determined
// based on the return value of this field.
// This field is only effective within the range of MaxIdemponentCallAttempts.
func WithRetryIfErr(retryIfErr fasthttp.RetryIfErrFunc) Option {
	return func(c *fasthttp.Client) {
		c.RetryIfErr = retryIfErr
	}
}

// ConfigureClient configures the fasthttp.HostClient.
func WithConfigureClient(configureClient func(hc *fasthttp.HostClient) error) Option {
	return func(c *fasthttp.Client) {
		c.ConfigureClient = configureClient
	}
}

// Client name. Used in User-Agent request header.
//
// Default client name is used if not set.
func WithName(name string) Option {
	return func(c *fasthttp.Client) {
		c.Name = name
	}
}

// Maximum number of connections per each host which may be established.
//
// DefaultMaxConnsPerHost is used if not set.
func WithMaxConnsPerHost(maxConnsPerHost int) Option {
	return func(c *fasthttp.Client) {
		c.MaxConnsPerHost = maxConnsPerHost
	}
}

// Idle keep-alive connections are closed after this duration.
//
// By default idle connections are closed
// after DefaultMaxIdleConnDuration.
func WithMaxIdleConnDuration(maxIdleConnDuration time.Duration) Option {
	return func(c *fasthttp.Client) {
		c.MaxIdleConnDuration = maxIdleConnDuration
	}
}

// Keep-alive connections are closed after this duration.
//
// By default connection duration is unlimited.
func WithMaxConnDuration(maxConnDuration time.Duration) Option {
	return func(c *fasthttp.Client) {
		c.MaxConnDuration = maxConnDuration
	}
}

// Maximum number of attempts for idempotent calls.
//
// DefaultMaxIdemponentCallAttempts is used if not set.
func WithMaxIdemponentCallAttempts(maxIdemponentCallAttempts int) Option {
	return func(c *fasthttp.Client) {
		c.MaxIdemponentCallAttempts = maxIdemponentCallAttempts
	}
}

// Per-connection buffer size for responses' reading.
// This also limits the maximum header size.
//
// Default buffer size is used if 0.
func WithReadBufferSize(readBufferSize int) Option {
	return func(c *fasthttp.Client) {
		c.ReadBufferSize = readBufferSize
	}
}

// Per-connection buffer size for requests' writing.
//
// Default buffer size is used if 0.
func WithWriteBufferSize(writeBufferSize int) Option {
	return func(c *fasthttp.Client) {
		c.WriteBufferSize = writeBufferSize
	}
}

// Maximum duration for full response reading (including body).
//
// By default response read timeout is unlimited.
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(c *fasthttp.Client) {
		c.ReadTimeout = readTimeout
	}
}

// Maximum duration for full request writing (including body).
//
// By default request write timeout is unlimited.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(c *fasthttp.Client) {
		c.WriteTimeout = writeTimeout
	}
}

// Maximum response body size.
//
// The client returns ErrBodyTooLarge if this limit is greater than 0
// and response body is greater than the limit.
//
// By default response body size is unlimited.
func WithMaxResponseBodySize(maxResponseBodySize int) Option {
	return func(c *fasthttp.Client) {
		c.MaxResponseBodySize = maxResponseBodySize
	}
}

// Maximum duration for waiting for a free connection.
//
// By default will not waiting, return ErrNoFreeConns immediately.
func WithMaxConnWaitTimeout(maxConnWaitTimeout time.Duration) Option {
	return func(c *fasthttp.Client) {
		c.MaxConnWaitTimeout = maxConnWaitTimeout
	}
}

// Connection pool strategy. Can be either LIFO or FIFO (default).
func WithConnPoolStrategy(connPoolStrategy fasthttp.ConnPoolStrategyType) Option {
	return func(c *fasthttp.Client) {
		c.ConnPoolStrategy = connPoolStrategy
	}
}

// NoDefaultUserAgentHeader when set to true, causes the default
// User-Agent header to be excluded from the Request.
func WithNoDefaultUserAgentHeader(noDefaultUserAgentHeader bool) Option {
	return func(c *fasthttp.Client) {
		c.NoDefaultUserAgentHeader = noDefaultUserAgentHeader
	}
}

// Attempt to connect to both ipv4 and ipv6 addresses if set to true.
//
// This option is used only if default TCP dialer is used,
// i.e. if Dial is blank.
//
// By default client connects only to ipv4 addresses,
// since unfortunately ipv6 remains broken in many networks worldwide :)
func WithDialDualStack(dialDualStack bool) Option {
	return func(c *fasthttp.Client) {
		c.DialDualStack = dialDualStack
	}
}

// Header names are passed as-is without normalization
// if this option is set.
//
// Disabled header names' normalization may be useful only for proxying
// responses to other clients expecting case-sensitive
// header names. See https://github.com/valyala/fasthttp/issues/57
// for details.
//
// By default request and response header names are normalized, i.e.
// The first letter and the first letters following dashes
// are uppercased, while all the other letters are lowercased.
// Examples:
//
//   - HOST -> Host
//   - content-type -> Content-Type
//   - cONTENT-lenGTH -> Content-Length
func WithDisableHeaderNamesNormalizing(disableHeaderNamesNormalizing bool) Option {
	return func(c *fasthttp.Client) {
		c.DisableHeaderNamesNormalizing = disableHeaderNamesNormalizing
	}
}

// Path values are sent as-is without normalization.
//
// Disabled path normalization may be useful for proxying incoming requests
// to servers that are expecting paths to be forwarded as-is.
//
// By default path values are normalized, i.e.
// extra slashes are removed, special characters are encoded.
func WithDisablePathNormalizing(disablePathNormalizing bool) Option {
	return func(c *fasthttp.Client) {
		c.DisablePathNormalizing = disablePathNormalizing
	}
}

// StreamResponseBody enables response body streaming.
func WithStreamResponseBody(streamResponseBody bool) Option {
	return func(c *fasthttp.Client) {
		c.StreamResponseBody = streamResponseBody
	}
}
