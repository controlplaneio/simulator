package cli

import (
	"context"
	"fmt"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			stateBucket, stateKey, name, err := getTerraformFlags(cmd)
			if err != nil {
				return fmt.Errorf("unable to get terraform flags: %w", err)
			}

			err = manager.Create(ctx, stateBucket, stateKey, name)
			if err != nil {
				return fmt.Errorf("unable to create simulator infrastructure: %w", err)
			}
			return nil
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			stateBucket, stateKey, name, err := getTerraformFlags(cmd)
			if err != nil {
				return fmt.Errorf("unable to get terraform flags: %w", err)
			}

			err = manager.Destroy(ctx, stateBucket, stateKey, name)
			if err != nil {
				return fmt.Errorf("unable to destroy simulator infrastructure: %w", err)
			}
			return nil
		},
	}

	for _, opt := range opts {
		opt(infraDestroyCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(infraDestroyCmd)
	}
}

func getTerraformFlags(cmd *cobra.Command) (string, string, string, error) {
	stateBucket, err := cmd.Flags().GetString("stateBucket")
	if err != nil {
		return "", "", "", fmt.Errorf("unable to get stateBucket flag: %w", err)
	}

	stateKey, err := cmd.Flags().GetString("stateKey")
	if err != nil {
		return "", "", "", fmt.Errorf("unable to get stateKey flag: %w", err)
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", "", "", fmt.Errorf("unable to get name flag: %w", err)
	}

	return stateBucket, stateKey, name, nil
}

func WithFlag(name, value, usage string) SimulatorCmdOptions {
	return func(command *cobra.Command) {
		command.Flags().String(name, value, usage)
	}
}
