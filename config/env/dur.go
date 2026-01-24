package env // import "go.microcore.dev/framework/config/env"

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.microcore.dev/framework/shutdown"
)

func MustDur(key string) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required duration variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		logger.Error(
			"failed to parse duration value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return d
}

func Dur(key string) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, fmt.Errorf("variable %s is not set", key)
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s duration value: %w", key, err)
	}

	return d, nil
}

func DurDefault(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		logger.Warn(
			"failed to parse duration value, using default",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Duration("default", def),
		)
		return def
	}

	return d
}
