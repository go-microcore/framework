package env // import "go.microcore.dev/framework/config/env"

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"go.microcore.dev/framework/shutdown"
)

func MustInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required int variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		logger.Error(
			"failed to parse int value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return i
}

func Int(key string) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, fmt.Errorf("variable %s is not set", key)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s int value: %w", key, err)
	}

	return i, nil
}

func IntDefault(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		logger.Warn(
			"failed to parse int value, using default",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Int("default", def),
		)
		return def
	}

	return i
}
