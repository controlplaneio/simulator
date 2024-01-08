package logging

import (
	"fmt"
	"log/slog"
	"os"
)

func Configure(level string) error {
	var leveler slog.Leveler

	switch level {
	case "error":
		leveler = slog.LevelError
	case "warn":
		leveler = slog.LevelWarn
	case "info":
		leveler = slog.LevelInfo
	case "debug":
		leveler = slog.LevelDebug
	default:
		return fmt.Errorf("unsupported log level: " + level)
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     leveler,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	slog.SetDefault(logger)

	return nil
}
