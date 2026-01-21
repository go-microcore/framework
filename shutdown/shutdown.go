package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

type (
	shutdown struct {
		signals  []os.Signal
		done     chan string
		manual   chan string
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

	Handler func(ctx context.Context, reason string) error
)

var (
	s = &shutdown{
		signals:  shutdownSignals,
		done:     make(chan string, 1),
		manual:   make(chan string, 1),
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
		return nil, errors.New("shutdown context has already been initialized")
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
	logger.Debug(
		"handler added",
		slog.Int("total", len(s.handlers)),
	)
}

func Wait() string {
	return <-s.done
}

func Shutdown(reason string) {
	s.manual <- reason
}

func Exit(reason string) {
	Shutdown(reason)
	logger.Info(
		"exit",
		slog.String("reason", Wait()),
	)
}

func subscribe() {
	logger.Debug("subscribe")

	signal.Notify(s.catch, s.signals...)
	defer signal.Stop(s.catch)

	var reason string

	select {
	case reason = <-s.manual:
	case sig := <-s.catch:
		reason = fmt.Sprintf("syscall (%s)", sig.String())
	}

	logger.Info(
		"shutdown",
		slog.String("reason", reason),
	)

	if s.ctx.cancel != nil {
		s.ctx.cancel()
	}

	handlers(reason)
	done(reason)
}

func handlers(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(s.handlers))

	for _, fn := range s.handlers {
		go func(fn Handler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						fmt.Sprintf("panic in shutdown handler: %v", r),
					)
				}
			}()
			if err := fn(ctx, reason); err != nil {
				logger.Error(
					"handler error",
					slog.Any("error", err),
				)
			}
		}(fn)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("shutdown handlers timed out")
	case <-done:
		logger.Debug("all shutdown handlers completed")
	}
}

func done(reason string) {
	s.once.Do(func() {
		s.done <- reason
		close(s.done)
		close(s.manual)
	})
}
