package env // import "go.microcore.dev/framework/config/env"

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"go.microcore.dev/framework/shutdown"
)

func MustBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required bool variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		logger.Error(
			"failed to parse bool value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return b
}

func Bool(key string) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return false, fmt.Errorf("variable %s is not set", key)
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s bool value: %w", key, err)
	}

	return b, nil
}

func BoolDefault(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		logger.Warn(
			"failed to parse bool value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		return def
	}

	return b
}
