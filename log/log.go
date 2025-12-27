package log // import "go.microcore.dev/framework/log"

import (
	"log/slog"
	"os"

	_ "go.microcore.dev/framework"
)

var (
	logLevel *slog.LevelVar
	logger   *slog.Logger
)

func init() {
	logLevel = &slog.LevelVar{}
	logLevel.Set(DefaultLogLevel)

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     logLevel,
					AddSource: false,
				},
			),
		),
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
