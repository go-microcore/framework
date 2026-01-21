package env // import "go.microcore.dev/framework/config/env"

import (
	_ "go.microcore.dev/framework"

	"github.com/joho/godotenv"
)

func New(filenames ...string) error {
	return godotenv.Load(filenames...)
}
