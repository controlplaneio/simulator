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
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.CreateBucket(ctx, bucket)
		cobra.CheckErr(err)
	},
}

func init() {
	createBucketCmd.Flags().StringVar(&bucket, "name", "", "the name of the bucket to create")

	bucketCmd.AddCommand(createBucketCmd)
	simulatorCmd.AddCommand(bucketCmd)
}
