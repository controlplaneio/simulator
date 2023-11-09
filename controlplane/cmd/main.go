package main

import (
	"github.com/controlplaneio/simulator/controlplane/cli"
	"github.com/controlplaneio/simulator/logging"
)

func main() {
	logging.Configure()
	cli.Execute()
}
