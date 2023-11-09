package main

import (
	"github.com/controlplaneio/simulator/internal/cli"
	"github.com/controlplaneio/simulator/logging"
)

func main() {
	logging.Configure()
	cli.Execute()
}
