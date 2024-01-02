package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/aws"
	"github.com/controlplaneio/simulator/v2/internal/config"
)

func WithBucketCmd(opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	bucketCmd := &cobra.Command{
		Use:   "bucket",
		Short: "Manage the bucket used to store the Terraform state",
	}

	for _, opt := range opts {
		opt(bucketCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(bucketCmd)
	}
}

func WithCreateBucketCmd(config config.Config, manager aws.BucketManager) SimulatorCmdOptions {
	bucketCreateCommand := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if config.Bucket == "" {
				slog.Error("Bucket name not configured, use the 'config' flag to set it")
				os.Exit(1)
			}

			err := manager.Create(ctx, config.Bucket)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(bucketCreateCommand)
	}
}

func WithDeleteBucketCmd(config config.Config, manager aws.BucketManager) SimulatorCmdOptions {
	bucketCreateCommand := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			err := manager.Delete(ctx, config.Bucket)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(bucketCreateCommand)
	}
}
