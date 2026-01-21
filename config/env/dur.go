package env // import "go.microcore.dev/framework/config/env"

import (
	"errors"
	"os"
	"time"
)

func MustDur(key string) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required duration variable " + key + " is not set")
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic("env: failed to parse duration " + key + ": " + err.Error())
	}
	return d
}

func Dur(key string) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, errors.New("env: variable " + key + " is not set")
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, err
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
		return def
	}
	return d
}
