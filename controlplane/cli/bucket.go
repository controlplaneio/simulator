package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/controlplane"
)

var bucketCmd = &cobra.Command{
	Use: "bucket",
}

var createBucketCmd = &cobra.Command{
	Use: "create",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.CreateBucket(ctx, bucket)
	},
}

func init() {
	createBucketCmd.Flags().StringVar(&bucket, "name", "", "the name of the bucket to create")

	bucketCmd.AddCommand(createBucketCmd)
	simulatorCmd.AddCommand(bucketCmd)
}
