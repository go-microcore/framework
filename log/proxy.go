package log // import "go.microcore.dev/framework/log"

import (
	"context"
	"log/slog"
)

type ProxyHandler struct {
	attrs []slog.Attr
	group string
}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

func (h *ProxyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return slog.Default().Handler().Enabled(ctx, level)
}

func (h *ProxyHandler) Handle(ctx context.Context, r slog.Record) error {
	if len(h.attrs) > 0 {
		r.AddAttrs(h.attrs...)
	}
	if h.group != "" {
		r.AddAttrs(slog.Group(h.group))
	}
	return slog.Default().Handler().Handle(ctx, r)
}

func (h *ProxyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ProxyHandler{
		attrs: append(append([]slog.Attr{}, h.attrs...), attrs...),
		group: h.group,
	}
}

func (h *ProxyHandler) WithGroup(name string) slog.Handler {
	return &ProxyHandler{
		attrs: h.attrs,
		group: name,
	}
}
