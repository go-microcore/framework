package env // import "go.microcore.dev/framework/config/env"

import (
	"errors"
	"os"
	"strconv"
)

func MustInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required int variable " + key + " is not set")
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic("env: failed to parse int " + key + ": " + err.Error())
	}
	return i
}

func Int(key string) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, errors.New("env: variable " + key + " is not set")
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
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
		return def
	}
	return i
}
