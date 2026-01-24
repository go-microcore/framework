package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"go.microcore.dev/framework/shutdown"
)

func MustBytesHex(key string) []byte {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required hex variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	b, err := hex.DecodeString(v)
	if err != nil {
		logger.Error(
			"failed to decode hex value",
			slog.String("key", key),
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return b
}

func BytesHex(key string) ([]byte, error) {
	v := os.Getenv(key)
	if v == "" {
		return nil, fmt.Errorf("variable %s is not set", key)
	}

	b, err := hex.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s hex value: %w", key, err)
	}

	return b, nil
}

func BytesHexDefault(key string, def []byte) []byte {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := hex.DecodeString(v)
	if err != nil {
		logger.Warn(
			"failed to parse hex value, using default",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Any("default", def),
		)
		return def
	}

	return b
}
