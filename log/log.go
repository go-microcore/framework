package log // import "go.microcore.dev/framework/log"

/*
Package log provides a global, centralized logging system for Go applications
built on top of the standard slog package.

This package is designed to make logging simple, consistent, and flexible
across your application or service, while minimizing boilerplate.

Key features:

1. **Single global logger**
   - All loggers created via log.New, log.With, or log.WithGroup share a
     central configuration.
   - A default logger is automatically initialized at startup with
     sensible defaults and developer-friendly output.

2. **Dynamic configuration**
   - Change output destination (stdout, file, custom writer), format
     (plain text, JSON, or pretty-colored), and attribute transformation
     at runtime via Config().
   - Changes take effect immediately across all loggers.

3. **Flexible backends**
   - Internally uses a proxy handler to delegate log events.
   - Updating the backend updates all existing loggers automatically.
   - Always use the package's helper functions to ensure logs go through
     the dynamic backend.

4. **Structured and contextual logging**
   - Supports adding attributes and grouping fields for better organization.
   - Makes it easy to enrich logs with context or package-specific information.

5. **Global log level management**
   - The log level can be queried and updated at runtime.
   - Changes apply to all loggers created through the package.

6. **Concurrency-safe**
   - Logging operations are safe to use from multiple goroutines.
   - Backend updates are protected by a mutex; reading logs is lock-free.

7. **Convenient helpers**
   - Functions for common log levels: Debug, Info, Warn, Error,
     with or without context.
   - Use these instead of managing slog.Logger instances directly
     to maintain consistency.

Best practices:

- Always create loggers through log.New, log.With, or log.WithGroup.
- Use Config() to adjust format or output, and SetLevel / SetLevelStr
  to adjust verbosity.
- Prefer structured logging with attributes over string concatenation.

This approach ensures consistent, flexible logging across your application,
with runtime configurability and safe concurrent use.
*/

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/lmittmann/tint"
	"golang.org/x/term"

	_ "go.microcore.dev/framework"
)

type (
	Options struct {
		Writer      io.Writer
		Format      OutputFormat
		ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	}
)

var (
	level   *slog.LevelVar
	backend slog.Handler
	logger  *slog.Logger
	mu      sync.Mutex
)

func init() {
	SetDefaultState()
}

// SetDefaultState resets the global logger to its default configuration.
// Typically used at startup or in tests.
func SetDefaultState() {
	// Set default level
	level = &slog.LevelVar{}
	level.Set(DefaultLogLevel)

	// Set default backend
	Config(
		Options{
			Writer:      DefaultWriter,
			Format:      DefaultFormat,
			ReplaceAttr: DefaultPrettyReplaceAttr,
		},
	)

	// New proxy handler
	handler := NewProxyHandler()

	// Set default immutable logger
	logger = slog.New(handler)
}

// SetBackend replaces the current backend for the global logger.
// Usually, Config() is sufficient.
func SetBackend(h slog.Handler) {
	mu.Lock()
	defer mu.Unlock()
	backend = h
}

// Config initializes and applies the global logger configuration.
//
// It sets up the global slog.Logger according to the provided Options,
// controlling:
//
//   - output destination (stdout, file, or any io.Writer);
//   - log format (plain text, JSON, or colorized "pretty");
//   - optional attribute transformation via ReplaceAttr.
//
// Behavior by format:
//
//   - FormatText
//     Writes human-readable logs using slog.TextHandler.
//     Suitable for local development or text log files.
//
//   - FormatJSON
//     Writes structured JSON logs using slog.JSONHandler.
//     Recommended for production environments and log aggregation systems.
//
//   - FormatPretty
//     Writes developer-friendly, colorized logs using tint.Handler.
//     Colors are automatically disabled if the output is not a terminal.
//
// ReplaceAttr:
//
//   - If not nil, it is applied to each attribute before output.
//     Can be used to mask sensitive data, rename fields, or modify messages.
//   - If nil, attributes are written as-is.
//
// Global effect:
//
//   - Config replaces the backend of the global logger.
//   - All loggers created via log.New, log.With, or log.WithGroup
//     immediately start using the new configuration.
//
// Example:
//
//   log.Config(log.Options{
//       Writer: os.Stdout,
//       Format: log.FormatPretty,
//   })
func Config(opts Options) error {
	var handler slog.Handler

	switch opts.Format {
	case FormatText:
		handler = slog.NewTextHandler(
			opts.Writer,
			&slog.HandlerOptions{
				Level:       level,
				ReplaceAttr: opts.ReplaceAttr,
			},
		)
	case FormatJSON:
		handler = slog.NewJSONHandler(
			opts.Writer,
			&slog.HandlerOptions{
				Level:       level,
				ReplaceAttr: opts.ReplaceAttr,
			},
		)
	case FormatPretty:
		handler = tint.NewHandler(
			opts.Writer,
			&tint.Options{
				Level:       level,
				TimeFormat:  DefaultTimeFormat,
				NoColor:     !isTerminal(opts.Writer),
				ReplaceAttr: opts.ReplaceAttr,
			})
	default:
		return fmt.Errorf("format %s not implemented", opts.Format)
	}

	SetBackend(handler)
	return nil
}

// SetLevel sets the global log level.
// Example: log.SetLevel(slog.LevelDebug)
func SetLevel(l slog.Level) {
	level.Set(l)
}

// SetLevelStr sets the global log level from a string.
// Example: log.SetLevelStr("INFO")
func SetLevelStr(l string) error {
	return level.UnmarshalText([]byte(l))
}

// GetLevel returns the current global log level.
func GetLevel() slog.Level {
	return level.Level()
}

// New creates a logger scoped to a specific package or component.
// Example: log.New("users").Info("user created")
func New(pkg string) *slog.Logger {
	return logger.With(
		slog.String("pkg", pkg),
	)
}

// With creates a logger with additional context fields.
// Example: log.With("user", "alice", "role", "admin").Info("login")
func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

// WithGroup creates a logger with a group for attributes.
// Example: log.WithGroup("db").Info("query executed")
func WithGroup(name string) *slog.Logger {
	return logger.WithGroup(name)
}

// Handler returns the current global slog.Handler.
// Useful for integration or testing.
func Handler() slog.Handler {
	return logger.Handler()
}

// Enabled checks whether logging at the given level is enabled.
func Enabled(ctx context.Context, level slog.Level) bool {
	return logger.Enabled(ctx, level)
}

// Log outputs a message at the specified level with optional arguments.
// Example: log.Log(ctx, slog.LevelInfo, "message", "key", 123)
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}

// LogAttrs outputs a message with structured slog attributes.
// Example: log.LogAttrs(ctx, slog.LevelInfo, "user login",
//          slog.String("user", "alice"), slog.Int("id", 42))
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}

// Debug logs a message at Debug level.
// Example: log.Debug("debug message", "key", 123)
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// DebugContext logs a Debug-level message with context.
// Example: log.DebugContext(ctx, "debug with ctx")
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

// Info logs a message at Info level.
// Example: log.Info("user logged in", "user", "alice")
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// InfoContext logs an Info-level message with context.
// Example: log.InfoContext(ctx, "operation completed")
func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

// Warn logs a message at Warn level.
// Example: log.Warn("disk space low")
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// WarnContext logs a Warn-level message with context.
// Example: log.WarnContext(ctx, "rate limit approaching")
func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

// Error logs a message at Error level.
// Example: log.Error("failed to save file", "file", filename)
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// ErrorContext logs an Error-level message with context.
// Example: log.ErrorContext(ctx, "request failed")
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

// isTerminal checks if the writer is a terminal.
// Used to automatically enable or disable colors.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
