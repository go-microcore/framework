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
	Wait()
	Shutdown(code int)
	Exit(code int)
	Recover()
	SetShutdownTimeout(t time.Duration)
}
