package env // import "go.microcore.dev/framework/config/env"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/joho/godotenv"
)

var logger = log.New(pkg)

func New(filenames ...string) error {
	return godotenv.Load(filenames...)
}
