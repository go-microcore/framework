package core // import "go.microcore.dev/framework/transport/http/client/core"

import (
	"time"

	"github.com/valyala/fasthttp"

	_ "go.microcore.dev/framework"
)

const (
	pkg = "go.microcore.dev/framework/transport/http/client/core"

	// Client name. Used in User-Agent request header.
	DefaultClientName = "microcore"

	// Maximum number of connections per each host which may be established.
	DefaultClientMaxConnsPerHost = 512

	// Idle keep-alive connections are closed after this duration.
	DefaultClientMaxIdleConnDuration = 10 * time.Second

	// Keep-alive connections are closed after this duration.
	DefaultClientMaxConnDuration = time.Duration(0)

	// Maximum number of attempts for idempotent calls.
	DefaultClientMaxIdemponentCallAttempts = 5

	// Per-connection buffer size for responses' reading.
	// This also limits the maximum header size.
	DefaultClientReadBufferSize = 4096

	// Per-connection buffer size for requests' writing.
	DefaultClientWriteBufferSize = 4096

	// Maximum duration for full response reading (including body).
	DefaultClientReadTimeout = 10 * time.Second

	// Maximum duration for full request writing (including body).
	DefaultClientWriteTimeout = 10 * time.Second

	// Header names are passed as-is without normalization
	// if this option is set.
	DefaultClientDisableHeaderNamesNormalizing = false

	// Path values are sent as-is without normalization.
	DefaultClientDisablePathNormalizing = false
)

var (
	// Increase DNS cache time to an hour instead of Default minute
	defaultClientDial = (&fasthttp.TCPDialer{
		Concurrency:      4096,
		DNSCacheDuration: time.Hour,
	}).Dial
)
