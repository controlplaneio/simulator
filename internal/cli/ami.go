package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/aws"
	"github.com/controlplaneio/simulator/v2/core/tools"
)

func WithAMICmd(opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	amiCmd := &cobra.Command{
		Use:   "ami",
		Short: "Manage the Simulator AMIs",
	}

	for _, opt := range opts {
		opt(amiCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(amiCmd)
	}
}

func WithAMIListCmd(manager aws.AMIManager) SimulatorCmdOptions {
	amiListCmd := &cobra.Command{
		Use:   "list",
		Short: "List simulator AMIs",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()

			amis, err := manager.List(ctx)
			if err != nil {
				return fmt.Errorf("unable to list simulator AMIs: %w", err)
			}

			table := tablewriter.NewWriter(os.Stdout)

			table.SetHeader([]string{
				"ID",
				"Type",
				"Name",
				"Created",
			})

			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)

			for _, ami := range amis {
				table.Append([]string{
					ami.ID,
					ami.Tags["Type"],
					ami.Name,
					ami.CreationDate(),
				})
				table.SetRowLine(true)
			}
			table.Render()
			return nil
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(amiListCmd)
	}
}

func WithAMIDeleteCmd(manager aws.AMIManager) SimulatorCmdOptions {
	imageDeleteCommand := &cobra.Command{
		Use:   "delete [ami id]",
		Short: "Delete a simulator AMI",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			id := args[0]

			err := manager.Delete(ctx, id)
			if err != nil {
				return fmt.Errorf("unable to delete simulator AMI: %w", err)
			}
			return nil
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(imageDeleteCommand)
	}
}

// TODO: Add flags for containerd, runc, cni, and kubernetes version

func WithAmiBuildCmd(builder tools.AMIBuilder) SimulatorCmdOptions {
	imageBuildCmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Build the packer image",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			id := args[0]
			err := builder.Build(ctx, id)
			if err != nil {
				return fmt.Errorf("unable to build packer image: %w", err)
			}
			return nil
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(imageBuildCmd)
	}
}
