package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestDefaultState(t *testing.T) {
	defer SetDefaultState()

	t.Run("nil level", func(t *testing.T) {
		if level == nil {
			t.Fatal("level should not be nil")
		}
	})

	t.Run("nil backend", func(t *testing.T) {
		if backend == nil {
			t.Fatal("backend should not be nil")
		}
	})

	t.Run("nil logger", func(t *testing.T) {
		if logger == nil {
			t.Fatal("logger should not be nil")
		}
	})

	t.Run("proxy handler", func(t *testing.T) {
		if _, ok := logger.Handler().(*ProxyHandler); !ok {
			t.Fatal("logger handler is not proxy handler")
		}
	})

	t.Run("default level", func(t *testing.T) {
		if got := level.Level(); got != DefaultLogLevel {
			t.Fatalf("expected default level %v, got %v", DefaultLogLevel, got)
		}
	})
}

// Note: this test is not safe for concurrent logging; do not access `records` from multiple goroutines.
func TestSetBackend(t *testing.T) {
	defer SetDefaultState()

	var records []slog.Record

	mockBackend := &mockBackend{
		handleFn: func(ctx context.Context, r slog.Record) error {
			records = append(records, r)
			return nil
		},
		enabledFn: func(ctx context.Context, l slog.Level) bool {
			return true
		},
	}

	SetBackend(mockBackend)
	Info("test message")

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	if got := records[0].Message; got != "test message" {
		t.Errorf("got message %q, want %q", got, "test message")
	}
}

func TestConfig_AllFormats(t *testing.T) {
	defer SetDefaultState()

	formats := []OutputFormat{FormatText, FormatJSON, FormatPretty}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			var buf bytes.Buffer
			replaceCalled := false

			err := Config(Options{
				Writer: &buf,
				Format: format,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					replaceCalled = true
					if a.Key == slog.MessageKey {
						return slog.String(a.Key, "REPLACED")
					}
					return a
				},
			})
			if err != nil {
				t.Fatalf("Config returned error: %v", err)
			}

			Info("hello", "key", 31337)

			output := buf.String()
			if output == "" {
				t.Fatal("expected log output, got empty buffer")
			}

			if !replaceCalled {
				t.Fatal("ReplaceAttr callback was not called")
			}

			switch format {
			case FormatJSON:
				var rec map[string]any
				if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
					t.Fatalf("invalid JSON output: %v", err)
				}
				if rec["msg"] != "REPLACED" {
					t.Fatalf("expected msg=REPLACED, got %v", rec["msg"])
				}
				if rec["key"] != float64(31337) { // json.Unmarshal превращает int в float64
					t.Fatalf("expected key=31337, got %v", rec["key"])
				}
			case FormatText, FormatPretty:
				if !strings.Contains(output, "REPLACED") {
					t.Fatalf("log output does not contain replaced message, got: %s", output)
				}
				if !strings.Contains(output, "key=31337") {
					t.Fatalf("log output does not contain attribute, got: %s", output)
				}
			}
		})
	}
}

func TestConfig_UnsupportedFormat(t *testing.T) {
	defer SetDefaultState()

	format := "xml"

	err := Config(Options{
		Writer: &bytes.Buffer{},
		Format: OutputFormat(format),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != fmt.Sprintf("format %s not implemented", format) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfig_ReplaceAttrNil(t *testing.T) {
	defer SetDefaultState()
	var buf bytes.Buffer

	if err := Config(Options{
		Writer: &buf,
		Format: FormatText,
	}); err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	Info("test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Fatal("expected log message without ReplaceAttr")
	}
}

func TestSetLevel(t *testing.T) {
	defer SetDefaultState()

	want := slog.LevelDebug
	SetLevel(want)

	if got := level.Level(); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSetLevelStr(t *testing.T) {
	defer SetDefaultState()

	want := slog.LevelWarn.String()
	SetLevelStr(want)

	if got := level.Level().String(); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSetLevelStr_Invalid(t *testing.T) {
	defer SetDefaultState()
	err := SetLevelStr("INVALID_LEVEL")
	if err == nil {
		t.Fatal("expected error for invalid level, got nil")
	}
}

func TestGetLevel(t *testing.T) {
	defer SetDefaultState()

	want := DefaultLogLevel

	if got := GetLevel(); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestNew(t *testing.T) {
	defer SetDefaultState()

	var buf bytes.Buffer
	Config(Options{
		Writer: &buf,
		Format: FormatText,
	})

	New("test").Info("msg")

	want := "level=INFO msg=msg pkg=test"

	if !strings.Contains(buf.String(), want) {
		t.Errorf("expected in log output: %q not contains %q", buf.String(), want)
	}
}

func TestWith(t *testing.T) {
	defer SetDefaultState()

	var buf bytes.Buffer
	if err := Config(Options{
		Writer: &buf,
		Format: FormatText,
	}); err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	With("user", "alice", "role", "admin").Info("hello")

	out := buf.String()
	want := "level=INFO msg=hello user=alice role=admin"

	fmt.Println(out)

	if !strings.Contains(out, want) {
		t.Errorf("expected in log output: %q not contains %q", out, want)
	}
}

func TestWithGroup(t *testing.T) {
	defer SetDefaultState()

	var buf bytes.Buffer
	Config(Options{
		Writer: &buf,
		Format: FormatText,
	})

	WithGroup("mygroup").Info("msg", "user_id", 31337)

	out := buf.String()
	want := "level=INFO msg=msg mygroup.user_id=31337"

	if !strings.Contains(out, want) {
		t.Errorf("expected in log output: %q not contains %q", out, want)
	}
}

func TestHandler(t *testing.T) {
	defer SetDefaultState()

	if _, ok := Handler().(*ProxyHandler); !ok {
		t.Fatal("handler is not proxy handler")
	}
}

func TestEnabled(t *testing.T) {
	defer SetDefaultState()

	SetLevel(slog.LevelWarn)

	if !Enabled(context.Background(), slog.LevelError) {
		t.Fatal("expected Enabled true for Error level")
	}
	if Enabled(context.Background(), slog.LevelInfo) {
		t.Fatal("expected Enabled false for Info level")
	}
}

func TestLog(t *testing.T) {
	defer SetDefaultState()

	var buf bytes.Buffer
	if err := Config(Options{
		Writer: &buf,
		Format: FormatText,
	}); err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	Log(context.Background(), slog.LevelInfo, "hello", "user", "alice", "role", "admin")

	out := buf.String()
	want := "level=INFO msg=hello user=alice role=admin"

	if !strings.Contains(out, want) {
		t.Errorf("expected in log output: %q not contains %q", out, want)
	}
}

func TestLogAttrs(t *testing.T) {
	defer SetDefaultState()

	var buf bytes.Buffer
	if err := Config(Options{
		Writer: &buf,
		Format: FormatText,
	}); err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	LogAttrs(context.Background(), slog.LevelInfo, "user login",
		slog.String("user", "alice"),
		slog.String("role", "admin"),
	)

	out := buf.String()
	want := "level=INFO msg=\"user login\" user=alice role=admin"

	if !strings.Contains(out, want) {
		t.Errorf("expected in log output: %q not contains %q", out, want)
	}
}

func TestLevelHelpers(t *testing.T) {
	defer SetDefaultState()

	SetLevel(slog.LevelDebug)

	var buf bytes.Buffer
	if err := Config(Options{
		Writer: &buf,
		Format: FormatText,
	}); err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func()
		want string
	}{
		{"Debug", func() { Debug("debug msg", "key", 1) }, "debug msg"},
		{"DebugContext", func() { DebugContext(ctx, "debug ctx msg", "key", 2) }, "debug ctx msg"},
		{"Info", func() { Info("info msg", "key", 3) }, "info msg"},
		{"InfoContext", func() { InfoContext(ctx, "info ctx msg", "key", 4) }, "info ctx msg"},
		{"Warn", func() { Warn("warn msg", "key", 5) }, "warn msg"},
		{"WarnContext", func() { WarnContext(ctx, "warn ctx msg", "key", 6) }, "warn ctx msg"},
		{"Error", func() { Error("error msg", "key", 7) }, "error msg"},
		{"ErrorContext", func() { ErrorContext(ctx, "error ctx msg", "key", 8) }, "error ctx msg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.fn()

			out := buf.String()
			if !strings.Contains(out, tt.want) {
				t.Fatalf("log output does not contain message %q, got: %s", tt.want, out)
			}
		})
	}
}

func TestProxyHandler_EmptyGroupAndAttrs(t *testing.T) {
	defer SetDefaultState()

	h := NewProxyHandler()
	if got := h.WithAttrs(nil); got != h {
		t.Error("WithAttrs(nil) should return same handler")
	}

	if got := h.WithGroup(""); got != h {
		t.Error("WithGroup(\"\") should return same handler")
	}
}

type mockBackend struct {
	handleFn  func(context.Context, slog.Record) error
	enabledFn func(context.Context, slog.Level) bool
	withAttrs func([]slog.Attr) slog.Handler
	withGroup func(string) slog.Handler
}

func (m *mockBackend) Enabled(ctx context.Context, l slog.Level) bool {
	return m.enabledFn(ctx, l)
}

func (m *mockBackend) Handle(ctx context.Context, r slog.Record) error {
	return m.handleFn(ctx, r)
}

func (m *mockBackend) WithAttrs(attrs []slog.Attr) slog.Handler {
	if m.withAttrs != nil {
		return m.withAttrs(attrs)
	}
	return m
}

func (m *mockBackend) WithGroup(name string) slog.Handler {
	if m.withGroup != nil {
		return m.withGroup(name)
	}
	return m
}
