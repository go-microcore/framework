package env // import "go.microcore.dev/framework/config/env"

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"

	_ "go.microcore.dev/framework"

	"github.com/joho/godotenv"
)

func New(filenames ...string) error {
	return godotenv.Load(filenames...)
}

// STRING

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

// BOOL

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

// INT

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

// BYTES (base64)

func MustBytes(key string) []byte {
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

func Bytes(key string) ([]byte, error) {
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

func BytesDefault(key string, def []byte) []byte {
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

// DURATION

func MustDuration(key string) time.Duration {
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

func Duration(key string) (time.Duration, error) {
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

func DurationDefault(key string, def time.Duration) time.Duration {
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
