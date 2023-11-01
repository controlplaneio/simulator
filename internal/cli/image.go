package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/internal/container"
)

var imageCmd = &cobra.Command{
	Use: "image",
}

var template string

var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the packer image",
	Run: func(cmd *cobra.Command, args []string) {
		runner := container.New(cfg)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		command := []string{
			"image",
			"build",
			"--template",
			fmt.Sprintf("%s.pkr.hcl", template),
		}

		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

func init() {
	imageBuildCmd.Flags().StringVar(&template, "template", "", "the packer template to build; bastion, or k8s")
	imageCmd.AddCommand(imageBuildCmd)

	simulatorCmd.AddCommand(imageCmd)
}
