package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	"go.microcore.dev/framework/shutdown"
)

func MustBytesB64(key string) []byte {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required base64 variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		logger.Error(
			"failed to decode base64 value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return b
}

func BytesB64(key string) ([]byte, error) {
	v := os.Getenv(key)
	if v == "" {
		return nil, fmt.Errorf("variable %s is not set", key)
	}

	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s base64 value: %w", key, err)
	}

	return b, nil
}

func BytesB64Default(key string, def []byte) []byte {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		logger.Warn(
			"failed to parse base64 value, using default",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Any("default", def),
		)
		return def
	}

	return b
}
