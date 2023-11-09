package main

import (
	"github.com/controlplaneio/simulator/v2/internal/cli"
	"github.com/controlplaneio/simulator/v2/logging"
)

func main() {
	logging.Configure()
	cli.Execute()
}
