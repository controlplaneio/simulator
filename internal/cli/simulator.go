package cli

import (
	"fmt"
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

	simulator.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		logLevel, err := cmd.Flags().GetString("log-level222")
		if err != nil {
			return fmt.Errorf("unable to get log-level flag: %w", err)
		}

		err = logging.Configure(logLevel)
		if err != nil {
			return fmt.Errorf("unable to configure logging: %w", err)
		}

		return nil
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
