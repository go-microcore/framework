package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/joho/godotenv"
)

var logger = log.New(pkg)

func New(filenames ...string) error {
	return godotenv.Load(filenames...)
}

// bool

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
			"failed to parse bool value, using default",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Bool("default", def),
		)
		return def
	}

	return b
}

// int

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

// string

func Str(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("variable %s is not set", key)
	}

	return v, nil
}

func StrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}

// duration

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

// bytes (hex)

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

// bytes (base64)

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
