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
		runner := container.New(cfg)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		command := []string{
			"bucket",
			"create",
			"--name",
			cfg.Bucket,
		}

		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

func init() {
	bucketCmd.AddCommand(createBucketCmd)
	simulatorCmd.AddCommand(bucketCmd)
}
