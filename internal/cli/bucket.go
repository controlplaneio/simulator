package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/internal/container"
)

var bucketCmd = &cobra.Command{
	Use: "bucket",
}

var createBucketCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		command := []string{
			"bucket",
			"create",
			"--name",
			cfg.Bucket,
		}

		runner := container.New(cfg)
		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

func init() {
	bucketCmd.AddCommand(createBucketCmd)
	simulatorCmd.AddCommand(bucketCmd)
}
