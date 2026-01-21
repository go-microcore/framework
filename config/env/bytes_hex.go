package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/hex"
	"errors"
	"os"
)

func MustBytesHex(key string) []byte {
	v := os.Getenv(key)
	if v == "" {
		panic("env: required hex variable " + key + " is not set")
	}
	b, err := hex.DecodeString(v)
	if err != nil {
		panic("env: failed to decode hex " + key + ": " + err.Error())
	}
	return b
}

func BytesHex(key string) ([]byte, error) {
	v := os.Getenv(key)
	if v == "" {
		return nil, errors.New("env: variable " + key + " is not set")
	}
	b, err := hex.DecodeString(v)
	if err != nil {
		return nil, err
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
		return def
	}
	return b
}
