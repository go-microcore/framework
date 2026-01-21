package env // import "go.microcore.dev/framework/config/env"

import (
	"errors"
	"os"
)

func MustStr(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required variable " + key + " is not set")
	}
	return v
}

func Str(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", errors.New("env: variable " + key + " is not set")
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
