package cli

import (
	"os"

	"github.com/controlplaneio/simulator/v2/internal/logging"
	"github.com/spf13/cobra"
)

const (
	logLevelEnv     = "SIMULATOR_LOG_LEVEL"
	logLevelDefault = "error"
)

type SimulatorCmdOptions func(command *cobra.Command)

func NewSimulatorCmd(opts ...SimulatorCmdOptions) *cobra.Command {
	simulator := &cobra.Command{
		Use:   "simulator",
		Short: "Simulator CLI",
	}

	simulator.PersistentFlags().String("log-level", logLevel(), "Log level (error, warn, info, debug)")

	simulator.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		cobra.CheckErr(err)

		err = logging.Configure(logLevel)
		cobra.CheckErr(err)
	}

	for _, opt := range opts {
		opt(simulator)
	}

	return simulator
}

func logLevel() string {
	if l := os.Getenv(logLevelEnv); len(l) > 0 {
		return l
	}
	return logLevelDefault
}
