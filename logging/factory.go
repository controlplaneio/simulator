package logging

import (
	"log/slog"
	"os"
)

func Configure() {
	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	slog.SetDefault(logger)
}
