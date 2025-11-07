package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

var (
	s = &shutdown{
		signals:  shutdownSignals,
		done:     make(chan error, 1),
		manual:   make(chan error, 1),
		catch:    make(chan os.Signal, 2),
		handlers: []func(os.Signal) error{},
	}
	logger = log.New(pkg)
)

type shutdown struct {
	signals  []os.Signal
	done     chan error
	manual   chan error
	catch    chan os.Signal
	handlers []func(os.Signal) error
	once     sync.Once
	mu       sync.Mutex
}

func init() {
	go subscribe()
}

func AddHandler(handler func(os.Signal) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, handler)
	logger.Debug(
		"handler added",
		slog.Int("total", len(s.handlers)),
	)
}

func Wait() error {
    return <-s.done
}

func Shutdown(err error) {
	s.manual <- err
}

func subscribe() {
	logger.Info("subscribe")
	signal.Notify(s.catch, s.signals...)
	defer signal.Stop(s.catch)
	select {
	case err := <-s.manual:
		logger.Info(
			"manual",
			slog.String("reason", err.Error()),
		)
		handlers(nil)
		done(err)
		return
	case sig := <-s.catch:
		logger.Info(
			"signal",
			slog.String("reason", sig.String()),
		)
		handlers(sig)
		done(fmt.Errorf("syscall (%s)", sig))
		return
	}
}

func handlers(sig os.Signal) {
	var wg sync.WaitGroup
	wg.Add(len(s.handlers))
	for _, fn := range s.handlers {
		go func(fn func(os.Signal) error) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						fmt.Sprintf("panic in shutdown handler: %v", r),
					)
				}
			}()
			if err := fn(sig); err != nil {
				logger.Error(
					"handler error",
					slog.Any("error", err),
				)
			}
		}(fn)
	}
	wg.Wait()
	logger.Info("all shutdown handlers completed")
}

func done(err error) {
	s.once.Do(func() {
		s.done <- err
		close(s.done)
		close(s.manual)
	})
}
