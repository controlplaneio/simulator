package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/controlplane"
)

var infraCmd = &cobra.Command{
	Use: "infra",
}

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.CreateInfrastructure(ctx, bucket, key, name)
		cobra.CheckErr(err)
	},
}

var destroyCmd = &cobra.Command{
	Use: "destroy",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.DestroyInfrastructure(ctx, bucket, key, name)
		cobra.CheckErr(err)
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
