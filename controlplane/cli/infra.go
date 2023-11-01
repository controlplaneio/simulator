package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/controlplane"
)

var infraCmd = &cobra.Command{
	Use: "infra",
}

var createCmd = &cobra.Command{
	Use: "create",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.CreateInfrastructure(ctx, bucket, key, name)
	},
}

var destroyCmd = &cobra.Command{
	Use: "destroy",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.DestroyInfrastructure(ctx, bucket, key, name)
	},
}

func init() {
	infraCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "the s3 bucket to use")
	infraCmd.PersistentFlags().StringVar(&key, "key", "state/terraform.tfstate", "the key to store state in the s3 bucket")
	infraCmd.PersistentFlags().StringVar(&name, "name", "", "the name for the infrastructure")

	infraCmd.AddCommand(createCmd)
	infraCmd.AddCommand(destroyCmd)
	simulatorCmd.AddCommand(infraCmd)
}
