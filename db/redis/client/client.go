package client // import "go.microcore.dev/framework/db/redis/client"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/redis/go-redis/v9"
)

var logger = log.New(pkg)

func New(opts ...Option) *redis.Client {
	options := &redis.Options{
		Network:    defaultNetwork,
		Addr:       defaultAddr,
		ClientName: defaultClientName,
		DB:         defaultDb,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(options)

	logger.Debug(
		"client has been successfully created",
		slog.String("network", options.Network),
		slog.String("addr", options.Addr),
		slog.String("client_name", options.ClientName),
		slog.Any("dialer", options.Dialer),
		slog.Any("on_connect", options.OnConnect),
		slog.Int("protocol", options.Protocol),
		slog.String("username", options.Username),
		slog.String("password", options.Password),
		slog.Any("credentials_provider", options.CredentialsProvider),
		slog.Any("credentials_provider_context", options.CredentialsProviderContext),
		slog.Any("streaming_credentials_provider", options.StreamingCredentialsProvider),
		slog.Int("db", options.DB),
		slog.Int("max_retries", options.MaxRetries),
		slog.Duration("min_retry_backoff", options.MinRetryBackoff),
		slog.Duration("max_retry_backoff", options.MaxRetryBackoff),
		slog.Duration("dial_timeout", options.DialTimeout),
		slog.Duration("read_timeout", options.ReadTimeout),
		slog.Duration("write_timeout", options.WriteTimeout),
		slog.Bool("context_timeout_enabled", options.ContextTimeoutEnabled),
		slog.Int("read_buffer_size", options.ReadBufferSize),
		slog.Int("write_buffer_size", options.WriteBufferSize),
		slog.Bool("pool_fifo", options.PoolFIFO),
		slog.Int("pool_size", options.PoolSize),
		slog.Duration("pool_timeout", options.PoolTimeout),
		slog.Int("min_idle_conns", options.MinIdleConns),
		slog.Int("max_idle_conns", options.MaxIdleConns),
		slog.Int("max_active_conns", options.MaxActiveConns),
		slog.Duration("conn_max_idle_time", options.ConnMaxIdleTime),
		slog.Duration("conn_max_lifetime", options.ConnMaxLifetime),
		slog.Any("tls_config", options.TLSConfig),
		slog.Any("limiter", options.Limiter),
		slog.Bool("disable_identity", options.DisableIdentity),
		slog.String("identity_suffix", options.IdentitySuffix),
		slog.Bool("unstable_resp3", options.UnstableResp3),
		slog.Int("failing_timeout_seconds", options.FailingTimeoutSeconds),
	)

	return client
}
