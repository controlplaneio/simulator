package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `create`,
		Short: "Runs terraform to create the required infrastructure for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stderr, "Not Implemented")

			return nil
		},
	}

	return cmd
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `destroy`,
		Short: "Tears down the infrastructure created for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stderr, "Not Implemented")

			return nil
		},
	}

	return cmd
}

func newDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `destroy`,
		Short: "Tears down the infrastructure created for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stderr, "Not Implemented")

			return nil
		},
	}

	return cmd
}

func newInfraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `infra <subcommand>`,
		Short:         "Interact with AWS to create, query and destroy the required infrastructure for scenarios",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newDestroyCommand())

	return cmd
}
