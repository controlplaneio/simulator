package main

import (
	"log/slog"
	"os"

	"github.com/controlplaneio/simulator/controlplane/cli"
)

const (
	LogLevel = "LOG_LEVEL"
)

func main() {
	level, ok := os.LookupEnv(LogLevel)
	if !ok {
		level = "info"
	}

	var sLevel slog.Level

	switch level {
	case "debug":
		sLevel = slog.LevelDebug
	case "info":
		sLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: sLevel,
	}))
	slog.SetDefault(logger)

	cli.Execute()
}
