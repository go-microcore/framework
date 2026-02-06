package client // import "go.microcore.dev/framework/db/redis/client"

import (
	"context"
	"fmt"
	"log/slog"

	_ "go.microcore.dev/framework"
)

type redisLogger struct {
	log *slog.Logger
}

func (r redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	r.log.Debug(fmt.Sprintf(format, v...))
}

func NewRedisLogger(log *slog.Logger) *redisLogger {
	return &redisLogger{
		log: log,
	}
}
