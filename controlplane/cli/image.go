package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/controlplane"
)

var (
	template string
)

var imageCmd = &cobra.Command{
	Use: "image",
}

var buildCmd = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.BuildImage(ctx, template)
		cobra.CheckErr(err)
	},
}

func init() {
	buildCmd.Flags().StringVar(&template, "template", "", "the packer template to build")

	imageCmd.AddCommand(buildCmd)
	simulatorCmd.AddCommand(imageCmd)
}
