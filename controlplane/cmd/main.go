package main

import (
	"github.com/controlplaneio/simulator/v2/controlplane/cli"
	"github.com/controlplaneio/simulator/v2/logging"
)

func main() {
	logging.Configure()
	cli.Execute()
}
