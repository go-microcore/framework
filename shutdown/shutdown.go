package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"context"
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
	}
	
	logger = log.New(pkg)
)

func init() {
	go subscribe()
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
