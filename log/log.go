package log // import "go.microcore.dev/framework/log"

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/term"

	_ "go.microcore.dev/framework"
)

type OutputFormat string

const (
	// FormatText writes plain text logs, human-readable, suitable for local or file output.
	FormatText OutputFormat = "text"

	// FormatJSON writes structured JSON logs, recommended for production and log aggregation.
	FormatJSON OutputFormat = "json"

	// FormatPretty writes colorized, developer-friendly logs, automatically disables
	// color if output is not a terminal.
	FormatPretty OutputFormat = "pretty"
)

type Options struct {
	Writer      io.Writer
	Format      OutputFormat
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

var (
	logLevel *slog.LevelVar
	logger   *slog.Logger
)

func init() {
	logLevel = &slog.LevelVar{}
	logLevel.Set(DefaultLogLevel)

	Configure(
		Options{
			Writer:      DefaultWriter,
			Format:      DefaultFormat,
			ReplaceAttr: DefaultPrettyReplaceAttr,
		},
	)

	logger = slog.New(
		NewProxyHandler(),
	)
}

func New(pkg string) *slog.Logger {
	return logger.With(
		slog.String("pkg", pkg),
	)
}

func Level() slog.Level {
	return logLevel.Level()
}

func SetLevelStr(level string) error {
	return logLevel.UnmarshalText([]byte(level))
}

func SetLevel(level slog.Level) {
	logLevel.Set(level)
}

// Configure initializes and installs the global logger configuration.
//
// It sets the default slog logger according to the provided Options,
// allowing to control:
//
//   - output destination (stdout, file, custom writer);
//   - log format (plain text, JSON, or pretty colored output);
//   - attribute transformation via ReplaceAttr.
//
// Behavior by format:
//
//   - FormatText
//     Uses slog.TextHandler and writes human-readable logs.
//     Intended for simple local usage or file-based text logs.
//
//   - FormatJSON
//     Uses slog.JSONHandler and produces structured JSON logs.
//     Recommended for production environments and log aggregation systems.
//
//   - FormatPretty
//     Uses tint.Handler to produce colorized, developer-friendly output.
//     Colors are enabled automatically only when the writer is a terminal;
//     when the output is redirected (e.g. to a file), colors are disabled.
//
// ReplaceAttr:
//
//	If ReplaceAttr is non-nil, it is applied to every log attribute
//	before it is written. This can be used to:
//	  - mask sensitive data;
//	  - rename or drop attributes;
//	  - customize time, level, or message formatting.
//
//	If ReplaceAttr is nil, attributes are written as-is.
//
// Global effects:
//
//	Configure replaces the global slog default logger via slog.SetDefault.
//	All loggers created via log.New(...) will immediately start using
//	the new configuration.
//
// Thread-safety:
//
//	Configure is expected to be called during application startup.
//	Reconfiguring the logger at runtime is supported by slog but may
//	lead to mixed log formats in concurrent systems and is discouraged.
//
// Example:
//
//	log.Configure(log.Options{
//	    Writer: os.Stdout,
//	    Format: log.FormatPretty,
//	})
func Configure(opts Options) error {
	var handler slog.Handler

	switch opts.Format {
	case FormatText:
		handler = slog.NewTextHandler(
			opts.Writer,
			&slog.HandlerOptions{
				Level:       logLevel,
				ReplaceAttr: opts.ReplaceAttr,
			},
		)
	case FormatJSON:
		handler = slog.NewJSONHandler(
			opts.Writer,
			&slog.HandlerOptions{
				Level:       logLevel,
				ReplaceAttr: opts.ReplaceAttr,
			},
		)
	case FormatPretty:
		handler = tint.NewHandler(
			opts.Writer,
			&tint.Options{
				Level:       logLevel,
				TimeFormat:  DefaultTimeFormat,
				NoColor:     !isTerminal(opts.Writer),
				ReplaceAttr: opts.ReplaceAttr,
			})
	default:
		return fmt.Errorf("format %s not implemented", opts.Format)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
