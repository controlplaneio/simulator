package cli

import (
	"github.com/spf13/cobra"
)

type SimulatorCmdOptions func(command *cobra.Command)

func NewSimulatorCmd(opts ...SimulatorCmdOptions) *cobra.Command {
	simulator := &cobra.Command{
		Use:   "simulator",
		Short: "Simulator CLI",
	}

	for _, opt := range opts {
		opt(simulator)
	}

	return simulator
}
