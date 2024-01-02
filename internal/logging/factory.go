package logging

import (
	"log/slog"
	"os"
)

func Configure() {
	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	slog.SetDefault(logger)
}
