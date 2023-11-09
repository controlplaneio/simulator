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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.BuildImage(ctx, template)
	},
}

func init() {
	buildCmd.Flags().StringVar(&template, "template", "", "the packer template to build")

	imageCmd.AddCommand(buildCmd)
	simulatorCmd.AddCommand(imageCmd)
}
