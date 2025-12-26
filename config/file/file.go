package file // import "go.microcore.dev/framework/config/file"

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	_ "go.microcore.dev/framework"

	"gopkg.in/yaml.v2"
)

type (
	format string

	c struct {
		path   string
		format format
		out    any
	}
)

const (
	// Formats
	JSON format = "json"
	YAML format = "yaml"
)

func New(opts ...Option) error {
	c := &c{
		path:   DefaultPath,
		format: DefaultFormat,
		out:    nil,
	}

	for _, opt := range opts {
		opt(c)
	}

	v := reflect.ValueOf(c.out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("config: out must be a non-nil pointer")
	}

	data, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("config: failed to read file %q: %w", c.path, err)
	}

	switch c.format {
	case JSON:
		if err := json.Unmarshal(data, c.out); err != nil {
			return fmt.Errorf("config: failed to unmarshal json: %w", err)
		}
	case YAML:
		if err := yaml.Unmarshal(data, c.out); err != nil {
			return fmt.Errorf("config: failed to unmarshal yaml: %w", err)
		}
	default:
		return fmt.Errorf("config: unsupported format %q", c.format)
	}

	return nil
}
