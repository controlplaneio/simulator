package main

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/cli/cmd"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var (
	logger *zap.SugaredLogger
)

func main() {
	var err error

	// logger writes to stderr
	logger, err = NewLogger("info", "console")
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer logger.Sync()

	if err := cmd.Execute(); err != nil {
		e := err.Error()

		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}
