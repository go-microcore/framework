package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

type (
	shutdown struct {
		exit     chan int
		code     chan int
		catch    chan os.Signal
		handlers []Handler
		once     sync.Once
		mu       sync.Mutex
		ctx      ctx
	}

	ctx struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	Handler func(ctx context.Context, code int) error
)

var (
	s = &shutdown{
		exit:     make(chan int, 1),
		code:     make(chan int, 1),
		catch:    make(chan os.Signal, 1),
		handlers: []Handler{},
		ctx: ctx{
			ctx: context.Background(),
		},
	}

	logger = log.New(pkg)
)

func init() {
	go subscribe()
}

// NewContext creates a root, shutdown-aware context for the entire program.
//
// In any Go application, the root context serves as the base for all other
// derived contexts. It is typically used to propagate cancellation signals
// and deadlines across multiple goroutines and services. `NewContext` ensures
// that the root context is properly initialized and protected from multiple
// concurrent creations.
//
// This function wraps `WithContext(context.Background())` and returns a
// cancelable context that will be automatically canceled when the application
// initiates a shutdown, either manually via `Shutdown()`/`Exit()` or due to
// system signals (e.g., SIGINT, SIGTERM). Using this root context as the base
// for all other contexts ensures consistent handling of shutdown across the
// program.
//
// Example usage:
//
//	ctx, err := shutdown.NewContext()
//	if err != nil {
//	    panic(err)
//	}
//	// pass `ctx` to servers, workers, or any long-running routines
//
// By standardizing the creation of a root shutdown context, the application
// can gracefully cancel all operations and clean up resources consistently
// when shutting down.
func NewContext() (context.Context, error) {
	return WithContext(context.Background())
}

func WithContext(parent context.Context) (context.Context, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx.cancel != nil {
		return nil, errors.New("shutdown context already initialized")
	}
	if parent == nil {
		return nil, errors.New("parent context is nil")
	}
	s.ctx.ctx, s.ctx.cancel = context.WithCancel(parent)
	return s.ctx.ctx, nil
}

func Context() context.Context {
	return s.ctx.ctx
}

func AddHandler(handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, handler)
}

func Wait() {
	os.Exit(<-s.exit)
}

func Shutdown(code int) {
	s.code <- code
}

func Exit(code int) {
	Shutdown(code)
	Wait()
}

func Recover() {
	if r := recover(); r != nil {
		logger.Error(
			"panic",
			slog.Any("error", r),
			slog.String("stack", string(debug.Stack())),
		)
		Exit(ExitPanic)
	}
}

func subscribe() {
	signal.Notify(s.catch, signals...)
	defer signal.Stop(s.catch)

	var code int

	select {
	case code = <-s.code:
	case sig := <-s.catch:
		code = ExitSignalBase + int(sig.(syscall.Signal))
	}

	logger.Info(
		"shutdown",
		slog.Int("code", code),
	)

	if s.ctx.cancel != nil {
		s.ctx.cancel()
	}

	if handlers(code) {
		if code > ExitSignalBase {
			code = ExitOK
		}
	} else {
		if code > ExitSignalBase {
			code = ExitShutdownError
		}
	}

	exit(code)
}

func handlers(code int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(s.handlers))

	var success atomic.Bool
	success.Store(true)

	for _, fn := range s.handlers {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						"panic in handler",
						slog.Any("error", r),
						slog.String("stack", string(debug.Stack())),
					)
					success.Store(false)
				}
			}()
			if err := fn(ctx, code); err != nil {
				logger.Error(
					"error in handler",
					slog.Any("error", err),
				)
				success.Store(false)
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("handlers timed out")
		return false
	case <-done:
		logger.Debug("all handlers completed")
		return success.Load()
	}
}

func exit(code int) {
	s.once.Do(func() {
		logger.Info(
			"exit",
			slog.Int("code", code),
		)
		os.Stdout.Sync()
		os.Stderr.Sync()
		s.exit <- code
		close(s.exit)
		close(s.code)
	})
}
