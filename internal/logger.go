package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func InitLogger(level string) {
	logLevel := slog.LevelInfo

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	case "info":
		// Default level
	default:
		fmt.Fprintf(os.Stderr, "Unknown log level %q, defaulting to INFO\n", level)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)

				if i := strings.Index(src.File, "internal/"); i >= 0 {
					src.File = src.File[i:]
				} else if i := strings.Index(src.File, "cmd/"); i >= 0 {
					src.File = src.File[i:]
				} else {
					src.File = filepath.Base(src.File)
				}

				return slog.Any(slog.SourceKey, src)
			}

			return a
		},
	}))

	slog.SetDefault(logger)

	switch strings.ToLower(level) {
	case "", "debug", "info", "warn", "error":
		// Valid levels.
	default:
		slog.Warn("Unknown log level, using INFO", "provided", level)
	}
}
