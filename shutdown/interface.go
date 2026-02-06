package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"context"
	"time"
)

type Manager interface {
	NewContext() (context.Context, error)
	WithContext(parent context.Context) (context.Context, error)
	Context() context.Context
	AddHandler(handler Handler) error
	Wait() int
	Shutdown(code int)
	Exit(code int) int
	SetShutdownTimeout(t time.Duration)
}
