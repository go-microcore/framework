package log // import "go.microcore.dev/framework/log"

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	_ "go.microcore.dev/framework"
)

const (
	DefaultLogLevel = slog.LevelInfo
	DefaultFormat   = FormatPretty
)

var (
	DefaultWriter            = os.Stdout
	DefaultTimeFormat        = time.StampMilli
	DefaultPrettyReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "pkg" {
			return tint.Attr(243, a)
		}
		if a.Key == "status" {
			code := a.Value.Int64()
			switch {
			case code >= 200 && code < 300:
				return tint.Attr(2, a) // 2xx
			case code >= 400 && code < 500:
				return tint.Attr(3, a) // 4xx
			case code >= 500 && code < 600:
				return tint.Attr(1, a) // 5xx
			}
		}
		return a
	}
)
