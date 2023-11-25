package cli

import (
	"context"

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

func WithCreateBucketCmd(config config.Config, creator aws.BucketCreator) SimulatorCmdOptions {
	bucketCreateCommand := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			err := creator.Create(ctx, config.Bucket)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(bucketCreateCommand)
	}
}
