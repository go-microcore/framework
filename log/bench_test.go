package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func BenchmarkWithJsonProxyHandler(b *testing.B) {
	defer SetDefaultState()

	var buf bytes.Buffer
	Config(Options{
		Writer: &buf,
		Format: FormatJSON,
	})

	logger := With("module", "benchmark")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "hello", "key", i)
	}
}

func BenchmarkWithJsonHandler(b *testing.B) {
	defer SetDefaultState()

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(handler).With("module", "benchmark")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "hello", "key", i)
	}
}

func BenchmarkWithGroupJsonProxyHandler(b *testing.B) {
	defer SetDefaultState()

	var buf bytes.Buffer
	Config(Options{
		Writer: &buf,
		Format: FormatJSON,
	})

	logger := WithGroup("grp")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "hello", "key", i)
	}
}

func BenchmarkWithGroupJsonHandler(b *testing.B) {
	defer SetDefaultState()

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(handler).WithGroup("grp")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "hello", "key", i)
	}
}
