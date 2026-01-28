package log // import "go.microcore.dev/framework/log"

/*
ProxyHandler is a specialized Handler for the log package that allows centralized management
of the global logging backend while adding attributes and groups on the fly. The main purposes
of ProxyHandler are:

1. Flexibility: allows creating loggers with additional fields (WithAttrs) and groups (WithGroup)
   without directly modifying the main backend.
2. Centralized control: all loggers created through the package automatically use the current global
   backend, simplifying format and output configuration for the entire application.
3. Ease of use: developers can work with logs via convenient Debug, Info, Warn, Error functions and
   their context-aware variants, without worrying about Handler internals.

Features:

- Attributes and groups added via WithAttrs/WithGroup apply only to the specific logger and are
  wrapped in groups so that JSON or text logs reflect the proper attribute structure.
- Changing the global backend via Config/SetBackend is automatically seen by ProxyHandler.
- The Handler is safe for concurrent use, but unlike the standard slog.Handler, each Handle call
  involves processing attributes and groups, which has a slight impact on performance.

Benchmark results on Apple M1 (darwin/arm64):

| Benchmark                          | ns/op   | B/op  | allocs/op |
|------------------------------------|---------|-------|-----------|
| WithJsonProxyHandler               | 634.9   | 299   | 0         |
| WithJsonHandler (stdlib)           | 494.2   | 229   | 0         |
| WithGroupJsonProxyHandler          | 647.9   | 318   | 2         |
| WithGroupJsonHandler (stdlib)      | 536.2   | 340   | 0         |

Conclusion:

- ProxyHandler provides flexibility and convenience in managing the global backend and log structure.
- This flexibility comes with a minor performance cost and extra allocations, especially when using groups.
- For performance-critical paths, using a direct slog.Handler without ProxyHandler may be preferable.
*/

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

type ProxyHandler struct {
	backend *slog.Handler
	attrs   []slog.Attr
	groups  []string
}

var (
	attrSlicePool = sync.Pool{
		New: func() any {
			s := make([]slog.Attr, 0, 16)
			return &s
		},
	}
	anySlicePool = sync.Pool{
		New: func() any {
			s := make([]any, 0, 16)
			return &s
		},
	}
)

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		backend: &backend,
	}
}

func (h *ProxyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return (*h.backend).Enabled(ctx, level)
}

func (h *ProxyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Fast return if there are no attributes or groups
	if len(h.attrs) == 0 && len(h.groups) == 0 {
		return (*h.backend).Handle(ctx, r)
	}

	// Get a slice for attributes from the pool
	attrsPtr := attrSlicePool.Get().(*[]slog.Attr)
	attrs := *attrsPtr
	attrs = attrs[:0] // reset length

	// Append attributes from With(...)
	attrs = append(attrs, h.attrs...)

	// Append attributes from the Record
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	// Wrap attributes in groups from the end
	for i := len(h.groups) - 1; i >= 0; i-- {
		anyPtr := anySlicePool.Get().(*[]any)
		anyAttrs := *anyPtr
		anyAttrs = anyAttrs[:0]

		for _, a := range attrs {
			anyAttrs = append(anyAttrs, a)
		}

		attrs = attrs[:0]
		attrs = append(attrs, slog.Group(h.groups[i], anyAttrs...))

		// Return anyAttrs to the pool
		*anyPtr = anyAttrs
		anySlicePool.Put(anyPtr)
	}

	// Create a new Record with the grouped attributes
	nr := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	nr.AddAttrs(attrs...)

	// Return attrs slice to the pool
	*attrsPtr = attrs
	attrSlicePool.Put(attrsPtr)

	return (*h.backend).Handle(ctx, nr)
}

func (h *ProxyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &ProxyHandler{
		backend: h.backend,
		attrs:   append(slices.Clone(h.attrs), attrs...),
		groups:  h.groups,
	}
}

func (h *ProxyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &ProxyHandler{
		backend: h.backend,
		attrs:   h.attrs,
		groups:  append(slices.Clone(h.groups), name),
	}
}
