package log // import "go.microcore.dev/framework/log"

import (
	"context"
	"log/slog"
	"slices"
)

type proxyHandler struct {
	attrs []slog.Attr
	group string
}

func NewProxyHandler() *proxyHandler {
	return &proxyHandler{}
}

func (h *proxyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return slog.Default().Handler().Enabled(ctx, level)
}

func (h *proxyHandler) Handle(ctx context.Context, r slog.Record) error {
	if len(h.attrs) > 0 {
		r.AddAttrs(h.attrs...)
	}
	if h.group != "" {
		r.AddAttrs(slog.Group(h.group))
	}
	return slog.Default().Handler().Handle(ctx, r)
}

func (h *proxyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &proxyHandler{
		attrs: append(slices.Clone(h.attrs), attrs...),
		group: h.group,
	}
}

func (h *proxyHandler) WithGroup(name string) slog.Handler {
	return &proxyHandler{
		attrs: h.attrs,
		group: name,
	}
}
