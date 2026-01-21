package env // import "go.microcore.dev/framework/config/env"

import (
	"errors"
	"os"
	"strconv"
)

func MustBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required bool variable " + key + " is not set")
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic("env: failed to parse bool " + key + ": " + err.Error())
	}
	return b
}

func Bool(key string) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return false, errors.New("env: variable " + key + " is not set")
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
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
		return def
	}
	return b
}
