package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/tools"
)

func WithInfraCmd(opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	amiCmd := &cobra.Command{
		Use:   "infra [command]",
		Short: "Manage the Simulator infrastructure",
	}

	for _, opt := range opts {
		opt(amiCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(amiCmd)
	}
}

func WithInfraCreateCmd(manager tools.InfraManager, opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	infraCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create simulator infrastructure",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			stateBucket, stateKey, name := getTerraformFlags(cmd)
			err := manager.Create(ctx, stateBucket, stateKey, name)
			cobra.CheckErr(err)
		},
	}

	for _, opt := range opts {
		opt(infraCreateCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(infraCreateCmd)
	}
}

func WithInfraDestroyCmd(manager tools.InfraManager, opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	infraDestroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy simulator infrastructure",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			stateBucket, stateKey, name := getTerraformFlags(cmd)
			err := manager.Destroy(ctx, stateBucket, stateKey, name)
			cobra.CheckErr(err)
		},
	}

	for _, opt := range opts {
		opt(infraDestroyCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(infraDestroyCmd)
	}
}

func getTerraformFlags(cmd *cobra.Command) (string, string, string) {
	stateBucket, err := cmd.Flags().GetString("stateBucket")
	cobra.CheckErr(err)

	stateKey, err := cmd.Flags().GetString("stateKey")
	cobra.CheckErr(err)

	name, err := cmd.Flags().GetString("name")
	cobra.CheckErr(err)
	return stateBucket, stateKey, name
}

func WithFlag(name, value, usage string) SimulatorCmdOptions {
	return func(command *cobra.Command) {
		command.Flags().String(name, value, usage)
	}
}
