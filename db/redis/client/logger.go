package client // import "go.microcore.dev/framework/db/redis/client"

import (
	"context"
	"log/slog"

	_ "go.microcore.dev/framework"
)

// redisLogger implements the redis.Logging interface to handle internal Redis messages.
type redisLogger struct {
	log *slog.Logger
}

// Printf is a stub that satisfies the redis.Logging interface. 
// Internal logging is currently silenced to prevent duplication with application-level logs.
func (r redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	// r.log.Debug(fmt.Sprintf(format, v...))
}

// NewRedisLogger creates a new instance of redisLogger with the provided slog instance.
func NewRedisLogger(log *slog.Logger) *redisLogger {
	return &redisLogger{
		log: log,
	}
}
