package env // import "go.microcore.dev/framework/config/env"

import (
	"fmt"
	"log/slog"
	"os"

	"go.microcore.dev/framework/shutdown"
)

func MustStr(key string) string {
	v := os.Getenv(key)
	if v == "" {
		logger.Error(
			"required string variable is not set",
			slog.String("key", key),
		)
		shutdown.Exit(shutdown.ExitConfigError)
	}

	return v
}

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
