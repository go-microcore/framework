package client // import "go.microcore.dev/framework/db/redis/client"

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	_ "go.microcore.dev/framework"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/auth"
)

type Option func(*redis.Options)

// Network type, either tcp or unix.
//
// default: is tcp.
func WithNetwork(network string) Option {
	return func(r *redis.Options) {
		r.Network = network
	}
}

// Addr is the address formated as host:port
func WithAddr(addr string) Option {
	return func(r *redis.Options) {
		r.Addr = addr
	}
}

// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
func WithClientName(clientName string) Option {
	return func(r *redis.Options) {
		r.ClientName = clientName
	}
}

// Dialer creates new network connection and has priority over
// Network and Addr options.
func WithDialer(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(r *redis.Options) {
		r.Dialer = dialer
	}
}

// Hook that is called when new connection is established.
func WithOnConnect(onConnect func(ctx context.Context, cn *redis.Conn) error) Option {
	return func(r *redis.Options) {
		r.OnConnect = onConnect
	}
}

// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
//
// default: 3.
func WithProtocol(protocol int) Option {
	return func(r *redis.Options) {
		r.Protocol = protocol
	}
}

// Username is used to authenticate the current connection
// with one of the connections defined in the ACL list when connecting
// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
func WithUsername(username string) Option {
	return func(r *redis.Options) {
		r.Username = username
	}
}

// Password is an optional password. Must match the password specified in the
// `requirepass` server configuration option (if connecting to a Redis 5.0 instance, or lower),
// or the User Password when connecting to a Redis 6.0 instance, or greater,
// that is using the Redis ACL system.
func WithPassword(password string) Option {
	return func(r *redis.Options) {
		r.Password = password
	}
}

// CredentialsProvider allows the username and password to be updated
// before reconnecting. It should return the current username and password.
func WithCredentialsProvider(credentialsProvider func() (username string, password string)) Option {
	return func(r *redis.Options) {
		r.CredentialsProvider = credentialsProvider
	}
}

// CredentialsProviderContext is an enhanced parameter of CredentialsProvider,
// done to maintain API compatibility. In the future,
// there might be a merge between CredentialsProviderContext and CredentialsProvider.
// There will be a conflict between them; if CredentialsProviderContext exists, we will ignore CredentialsProvider.
func WithCredentialsProviderContext(credentialsProviderContext func(ctx context.Context) (username string, password string, err error)) Option {
	return func(r *redis.Options) {
		r.CredentialsProviderContext = credentialsProviderContext
	}
}

// StreamingCredentialsProvider is used to retrieve the credentials
// for the connection from an external source. Those credentials may change
// during the connection lifetime. This is useful for managed identity
// scenarios where the credentials are retrieved from an external source.
//
// Currently, this is a placeholder for the future implementation.
func WithStreamingCredentialsProvider(streamingCredentialsProvider auth.StreamingCredentialsProvider) Option {
	return func(r *redis.Options) {
		r.StreamingCredentialsProvider = streamingCredentialsProvider
	}
}

// DB is the database to be selected after connecting to the server.
func WithDB(db int) Option {
	return func(r *redis.Options) {
		r.DB = db
	}
}

// MaxRetries is the maximum number of retries before giving up.
// -1 (not 0) disables retries.
//
// default: 3 retries
func WithMaxRetries(maxRetries int) Option {
	return func(r *redis.Options) {
		r.MaxRetries = maxRetries
	}
}

// MinRetryBackoff is the minimum backoff between each retry.
// -1 disables backoff.
//
// default: 8 milliseconds
func WithMinRetryBackoff(minRetryBackoff time.Duration) Option {
	return func(r *redis.Options) {
		r.MinRetryBackoff = minRetryBackoff
	}
}

// MaxRetryBackoff is the maximum backoff between each retry.
// -1 disables backoff.
// default: 512 milliseconds;
func WithMaxRetryBackoff(maxRetryBackoff time.Duration) Option {
	return func(r *redis.Options) {
		r.MaxRetryBackoff = maxRetryBackoff
	}
}

// DialTimeout for establishing new connections.
//
// default: 5 seconds
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(r *redis.Options) {
		r.DialTimeout = dialTimeout
	}
}

// ReadTimeout for socket reads. If reached, commands will fail
// with a timeout instead of blocking. Supported values:
//
//   - `-1` - no timeout (block indefinitely).
//   - `-2` - disables SetReadDeadline calls completely.
//
// default: 3 seconds
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(r *redis.Options) {
		r.ReadTimeout = readTimeout
	}
}

// WriteTimeout for socket writes. If reached, commands will fail
// with a timeout instead of blocking.  Supported values:
//
//   - `-1` - no timeout (block indefinitely).
//   - `-2` - disables SetWriteDeadline calls completely.
//
// default: 3 seconds
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(r *redis.Options) {
		r.WriteTimeout = writeTimeout
	}
}

// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
func WithContextTimeoutEnabled(contextTimeoutEnabled bool) Option {
	return func(r *redis.Options) {
		r.ContextTimeoutEnabled = contextTimeoutEnabled
	}
}

// ReadBufferSize is the size of the bufio.Reader buffer for each connection.
// Larger buffers can improve performance for commands that return large responses.
// Smaller buffers can improve memory usage for larger pools.
//
// default: 32KiB (32768 bytes)
func WithReadBufferSize(readBufferSize int) Option {
	return func(r *redis.Options) {
		r.ReadBufferSize = readBufferSize
	}
}

// WriteBufferSize is the size of the bufio.Writer buffer for each connection.
// Larger buffers can improve performance for large pipelines and commands with many arguments.
// Smaller buffers can improve memory usage for larger pools.
//
// default: 32KiB (32768 bytes)
func WithWriteBufferSize(writeBufferSize int) Option {
	return func(r *redis.Options) {
		r.WriteBufferSize = writeBufferSize
	}
}

// PoolFIFO type of connection pool.
//
//   - true for FIFO pool
//   - false for LIFO pool.
//
// Note that FIFO has slightly higher overhead compared to LIFO,
// but it helps closing idle connections faster reducing the pool size.
func WithPoolFIFO(poolFIFO bool) Option {
	return func(r *redis.Options) {
		r.PoolFIFO = poolFIFO
	}
}

// PoolSize is the base number of socket connections.
// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
// you can limit it through MaxActiveConns
//
// default: 10 * runtime.GOMAXPROCS(0)
func WithPoolSize(poolSize int) Option {
	return func(r *redis.Options) {
		r.PoolSize = poolSize
	}
}

// PoolTimeout is the amount of time client waits for connection if all connections
// are busy before returning an error.
//
// default: ReadTimeout + 1 second
func WithPoolTimeout(poolTimeout time.Duration) Option {
	return func(r *redis.Options) {
		r.PoolTimeout = poolTimeout
	}
}

// MinIdleConns is the minimum number of idle connections which is useful when establishing
// new connection is slow. The idle connections are not closed by default.
//
// default: 0
func WithMinIdleConns(minIdleConns int) Option {
	return func(r *redis.Options) {
		r.MinIdleConns = minIdleConns
	}
}

// MaxIdleConns is the maximum number of idle connections.
// The idle connections are not closed by default.
//
// default: 0
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(r *redis.Options) {
		r.MaxIdleConns = maxIdleConns
	}
}

// MaxActiveConns is the maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
// If the pool is full, the next call to Get() will block until a connection is released.
func WithMaxActiveConns(maxActiveConns int) Option {
	return func(r *redis.Options) {
		r.MaxActiveConns = maxActiveConns
	}
}

// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
// Should be less than server's timeout.
//
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are not closed due to a connection's idle time.
// -1 disables idle timeout check.
//
// default: 30 minutes
func WithConnMaxIdleTime(connMaxIdleTime time.Duration) Option {
	return func(r *redis.Options) {
		r.ConnMaxIdleTime = connMaxIdleTime
	}
}

// ConnMaxLifetime is the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
// If <= 0, connections are not closed due to a connection's age.
//
// default: 0
func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(r *redis.Options) {
		r.ConnMaxLifetime = connMaxLifetime
	}
}

// TLSConfig to use. When set, TLS will be negotiated.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(r *redis.Options) {
		r.TLSConfig = tlsConfig
	}
}

// Limiter interface used to implement circuit breaker or rate limiter.
func WithLimiter(limiter redis.Limiter) Option {
	return func(r *redis.Options) {
		r.Limiter = limiter
	}
}

// DisableIdentity is used to disable CLIENT SETINFO command on connect.
//
// default: false
func WithDisableIdentity(disableIdentity bool) Option {
	return func(r *redis.Options) {
		r.DisableIdentity = disableIdentity
	}
}

// Add suffix to client name. Default is empty.
// IdentitySuffix - add suffix to client name.
func WithIdentitySuffix(identitySuffix string) Option {
	return func(r *redis.Options) {
		r.IdentitySuffix = identitySuffix
	}
}

// UnstableResp3 enables Unstable mode for Redis Search module with RESP3.
// When unstable mode is enabled, the client will use RESP3 protocol and only be able to use RawResult
func WithUnstableResp3(unstableResp3 bool) Option {
	return func(r *redis.Options) {
		r.UnstableResp3 = unstableResp3
	}
}

// FailingTimeoutSeconds is the timeout in seconds for marking a cluster node as failing.
// When a node is marked as failing, it will be avoided for this duration.
// Default is 15 seconds.
func WithFailingTimeoutSeconds(failingTimeoutSeconds int) Option {
	return func(r *redis.Options) {
		r.FailingTimeoutSeconds = failingTimeoutSeconds
	}
}
