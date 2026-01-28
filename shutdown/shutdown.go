package shutdown // import "go.microcore.dev/framework/shutdown"

/*
Package shutdown provides a global shutdown manager for Go applications,
allowing centralized handling of application termination, context cancellation,
and execution of user-defined shutdown handlers.

Lifecycle of the Manager:

1. stateInit

   - Initial state.
   - The manager has just been created or the global defaultManager is not yet initialized.
   - Allowed actions:
       • Set a custom default manager (for testing).
       • Initialize the default manager.
       • Add shutdown handlers.
       • Create a root context.
   - After the first call to Default() or SetDefaultManager(), the state transitions to stateRunning.

2. stateRunning

   - The manager is active and the application is running.
   - Allowed actions:
       • Add shutdown handlers (before shutdown starts).
       • Initialize or read the context.
   - Restrictions:
       • Replacing the default manager is not allowed.
       • Once shutdown is triggered, the state will move to stateShuttingDown.

3. stateShuttingDown

   - The manager has received a shutdown request via Shutdown()/Exit() or an OS signal.
   - What happens:
       • State changes to stateShuttingDown.
       • The root context is canceled, notifying all dependent goroutines.
       • Registered shutdown handlers are executed.
           - Any errors or panics in handlers are logged.
   - Restrictions:
       • Adding new handlers is not allowed.
       • Replacing the default manager is not allowed.

4. stateExited

   - Final state after shutdown completes.
   - Actions:
       • Log the shutdown completion and exit code.
       • Exit the application.
   - After this, any calls to the manager are invalid or will result in errors.

Key points:

- The root context can be safely accessed from any goroutine.
- Handlers added after shutdown has started are not allowed.
- SetDefaultManager() is intended for testing and must be called before the default manager is initialized.
*/

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

const (
	stateInit state = iota
	stateRunning
	stateShuttingDown
	stateExited
)

type (
	state   int
	manager struct {
		state    atomic.Int32
		exit     chan int
		code     chan int
		catch    chan os.Signal
		handlers []Handler
		once     sync.Once
		mu       sync.Mutex
		ctx      struct {
			ctx    atomic.Value // context.Context
			cancel atomic.Value // context.CancelFunc
		}
	}
	Handler func(ctx context.Context, code int) error
)

var (
	defaultManager Manager
	defaultState   atomic.Int32
	timeout        time.Duration
	exitFunc       func(code int)
	logger         *slog.Logger
	once           sync.Once
)

func init() {
	defaultState.Store(int32(stateInit))
	timeout = defaultShutdownTimeout
	exitFunc = os.Exit
	logger = log.New(pkg)
}

func newManager() Manager {
	m := &manager{
		exit:     make(chan int, 1),
		code:     make(chan int, 1),
		catch:    make(chan os.Signal, 1),
		handlers: []Handler{},
	}
	m.state.Store(int32(stateInit))
	go m.subscribe()
	return m
}

func (m *manager) NewContext() (context.Context, error) {
	return m.WithContext(context.Background())
}

func (m *manager) WithContext(parent context.Context) (context.Context, error) {
	if m.ctx.ctx.Load() != nil {
		return nil, ErrContextAlreadyInit
	}
	if parent == nil {
		return nil, ErrParentContextNil
	}

	ctx, cancel := context.WithCancel(parent)
	m.ctx.ctx.Store(ctx)
	m.ctx.cancel.Store(cancel)

	return ctx, nil
}

func (m *manager) Context() context.Context {
	v := m.ctx.ctx.Load()
	if v == nil {
		return context.Background()
	}
	return v.(context.Context)
}

func (m *manager) AddHandler(handler Handler) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.Load() != int32(stateInit) {
		return ErrCannotAddHandlerAfterShutdown
	}
	m.handlers = append(m.handlers, handler)
	return nil
}

func (m *manager) Wait() {
	exitFunc(<-m.exit)
}

func (m *manager) Shutdown(code int) {
	select {
	case m.code <- code:
	default:
		logger.Warn("code not sent: channel blocked")
	}
}

func (m *manager) Exit(code int) {
	m.Shutdown(code)
	m.Wait()
}

func (m *manager) Recover() {
	if r := recover(); r != nil {
		logger.Error(
			"panic",
			slog.Any("error", r),
			slog.String("stack", string(debug.Stack())),
		)
		m.Exit(ExitPanic)
	}
}

func (m *manager) SetShutdownTimeout(t time.Duration) {
	timeout = t
}

func (m *manager) subscribe() {
	signal.Notify(m.catch, signals...)
	defer signal.Stop(m.catch)

	var code int

	select {
	case code = <-m.code:
	case sig := <-m.catch:
		code = ExitSignalBase + int(sig.(syscall.Signal))
	}

	m.state.Store(int32(stateShuttingDown))

	logger.Info(
		"shutdown",
		slog.Int("code", code),
	)

	if c := m.ctx.cancel.Load(); c != nil {
		c.(context.CancelFunc)()
	}

	if !m.exec(code) {
		code = ExitShutdownError
	} else if code > ExitSignalBase {
		code = ExitOK
	}

	m.term(code)
}

func (m *manager) exec(code int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(m.handlers))

	var success atomic.Bool
	success.Store(true)

	for _, fn := range slices.Clone(m.handlers) {
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
		s := success.Load()
		if s {
			logger.Debug("all handlers completed without errors")
		} else {
			logger.Warn("all handlers completed with errors")
		}
		return s
	}
}

func (m *manager) term(code int) {
	m.once.Do(func() {
		m.state.Store(int32(stateExited))
		logger.Info(
			"exit",
			slog.Int("code", code),
		)
		os.Stdout.Sync()
		os.Stderr.Sync()
		m.exit <- code
	})
}

func def() Manager {
	once.Do(func() {
		defaultManager = newManager()
		defaultState.Store(int32(stateRunning))
	})
	return defaultManager
}

// Public API

// SetDefaultManager replaces the global default shutdown manager.
//
// This function is intended mainly for testing or special scenarios where
// you need to provide a custom shutdown manager instead of the default one.
//
// How it works:
//
// 1. It checks the current lifecycle state of the default manager:
//   - If the manager has not yet been initialized, the custom manager is set.
//   - If the manager is already running or shutdown has started, the call fails with an error.
//     2. After successfully setting the manager, it becomes the global default
//     for all subsequent shutdown operations.
//
// Constraints:
//
// - Must be called before any call to Default().
// - Cannot replace the manager after shutdown has started.
//
// Example:
//
//	err := shutdown.SetDefaultManager(customMgr)
//	if err != nil {
//	    log.Fatal("failed to set custom shutdown manager:", err)
//	}
func SetDefaultManager(m Manager) error {
	if !defaultState.CompareAndSwap(
		int32(stateInit),
		int32(stateRunning),
	) {
		switch state(defaultState.Load()) {
		case stateRunning:
			return ErrManagerAlreadyRunning
		case stateShuttingDown, stateExited:
			return ErrCannotCallAfterShutdown
		default:
			return ErrUnknownState
		}
	}
	defaultManager = m
	return nil
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
// Example:
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
	return def().NewContext()
}

func WithContext(parent context.Context) (context.Context, error) {
	return def().WithContext(parent)
}

// Context returns the root shutdown-aware context managed by this manager.
//
// This function provides access to the application's root context, which is automatically
// canceled when the application begins shutting down. Using this context allows
// goroutines, workers, and services to detect shutdown events and stop gracefully.
//
// How it works:
//
// 1. Checks if a shutdown context has been initialized:
//   - If yes, returns the stored context.
//   - If no, returns a default background context to ensure the caller always gets
//     a valid context.
//     2. The returned context can be used anywhere in the application where a cancellable
//     context is needed.
//
// Constraints:
//
// - The returned context is read-only; it should not be re-assigned or replaced.
// - The context may be canceled automatically during shutdown.
//
// Example:
//
//	select {
//	case <-shutdown.Context().Done():
//	    // Shutdown signal received, stop work gracefully
//	default:
//	    // Continue normal operation
//	}
func Context() context.Context {
	return def().Context()
}

// AddHandler registers a new shutdown handler.
//
// A handler is a function that will be executed when the application
// is shutting down, either due to a system signal or a call to Shutdown()/Exit().
//
// How it works:
//
// 1. The manager's state is checked:
//
//   - If shutdown has not started yet, the handler is added to the list.
//
//   - If shutdown is already in progress, adding a handler is not allowed and an error is returned.
//
//     2. All registered handlers will be called during shutdown.
//     They receive the root context and an exit code, allowing services to
//     gracefully stop, release resources, and handle any cleanup or errors.
//
// Constraints:
//
// - Handlers can only be added before shutdown starts.
// - Any attempts to add handlers after shutdown has begun will fail.
//
// Example:
//
//	err := shutdown.AddHandler(func(ctx context.Context, code int) error {
//	    log.Println("Closing database connection...")
//	    // perform cleanup
//	    return nil
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
func AddHandler(handler Handler) error {
	return def().AddHandler(handler)
}

// Wait blocks the current goroutine until a shutdown signal is received.
//
// This function keeps the application running until a termination signal is triggered.
// It waits for one of the following events:
//   - A call to Shutdown() or Exit() within the application code.
//   - Receiving an OS termination signal (SIGINT, SIGTERM, SIGQUIT).
//
// How it works:
//
//  1. The function blocks until a shutdown code is sent to the exit channel.
//  2. Once the code is received, it immediately calls os.Exit(code) to terminate the process.
//  3. Any registered shutdown handlers are executed before the exit code is sent,
//     as Shutdown/Exit triggers them beforehand.
//
// Example:
//
//	func main() {
//	    // Setup root context
//	    ctx, _ := shutdown.NewContext()
//
//	    // Add shutdown handlers
//	    shutdown.AddHandler(func(ctx context.Context, code int) error {
//	        log.Println("Closing database connections")
//	        return nil
//	    })
//
//	    // Wait for shutdown signal
//	    shutdown.Wait() // Blocks until Shutdown/Exit is called or SIGINT/SIGTERM/SIGQUIT is received
//	}
func Wait() {
	def().Wait()
}

// Shutdown initiates program termination with the specified exit code.
//
// This function signals the manager that the application should shut down.
// It **does not block** execution at the call site — shutdown proceeds in the background,
// including execution of all registered shutdown handlers (graceful shutdown).
// The exit code will be used when the program finally terminates.
//
// How it works:
//
// 1. Sends the exit code to the manager's internal channel.
// 2. The manager, in a separate goroutine, starts the graceful shutdown process:
//   - cancels the root context,
//   - runs all registered shutdown handlers,
//   - exits the program with the specified code.
//
// Constraints:
//
// - Handlers may take some time to complete — shutdown proceeds asynchronously.
// - If the channel is blocked, the exit code might not be sent, and a warning is logged.
//
// Example:
//
//	// Initiate shutdown with exit code 0 (shutdown.ExitOK)
//	shutdown.Shutdown(shutdown.ExitOK)
//	// Program execution continues; graceful shutdown runs in parallel
func Shutdown(code int) {
	def().Shutdown(code)
}

// Exit initiates application termination with the specified exit code.
//
// This function starts the application shutdown process:
//  1. Triggers a graceful shutdown via Shutdown(code), notifying all registered
//     handlers and cancelling the root context.
//  2. Blocks the current goroutine until all shutdown handlers have finished their work.
//  3. After all handlers complete, the application immediately exits with the given code.
//
// Important:
//
// - Exit does not return until the entire shutdown process is complete.
// - It is used for final application termination in controlled situations.
//
// Example:
//
//		// Terminate the application with exit code 0 (shutdown.ExitOK) after all shutdown
//	 // handlers finish
//		shutdown.Exit(shutdown.ExitOK)
func Exit(code int) {
	def().Exit(code)
}

// Recover intercepts panics that occur in the application, logs them, and initiates a
// graceful shutdown.
//
// This function is used to protect the application from uncontrolled panics.
// If a panic occurs in any goroutine or part of the code, calling Recover() allows you to:
//   - capture error information and the stack trace,
//   - gracefully shut down the application via Exit() with a specific error code.
//
// How it works:
//
// 1. Checks whether a panic has occurred in the current goroutine using the built-in recover() function.
// 2. If a panic is detected:
//   - Logs a message including the panic details and stack trace.
//   - Calls Exit() with a predefined panic exit code (ExitPanic) to terminate the application gracefully.
//
// Example:
//
//	func main() {
//	    defer shutdown.Recover()
//
//	    // main application logic
//	    runServer()
//	}
func Recover() {
	def().Recover()
}

// SetShutdownTimeout sets the maximum duration allowed for shutdown handlers to complete.
//
// By default, each registered shutdown handler has a fixed amount of time to finish its work
// before being considered timed out. This function allows you to override that default
// timeout for the current shutdown manager.
//
// How it works:
//
// 1. Sets the internal `timeout` used by the manager when executing shutdown handlers.
// 2. All subsequent calls to shutdown will use this new timeout value.
// 3. If handlers do not complete within the timeout, they are canceled, and a warning is logged.
//
// Example:
//
//	// Allow handlers up to 5 seconds to complete
//	shutdown.SetShutdownTimeout(5 * time.Second)
//
// Important:
//
// - This timeout affects all handlers added before and after the call.
// - Should be called before triggering Shutdown/Exit for predictable behavior.
func SetShutdownTimeout(t time.Duration) {
	def().SetShutdownTimeout(t)
}
