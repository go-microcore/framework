package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/base64"
	"errors"
	"os"
)

func MustBytesB64(key string) []byte {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required base64 variable " + key + " is not set")
	}
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		panic("env: failed to decode base64 " + key + ": " + err.Error())
	}
	return b
}

func BytesB64(key string) ([]byte, error) {
	v := os.Getenv(key)
	if v == "" {
		return nil, errors.New("env: variable " + key + " is not set")
	}
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, err
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
		return def
	}
	return b
}
