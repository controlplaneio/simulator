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

// TODO: Add flags for containerd, runc, cni, and kubernetes version
// TODO: Add image cleanup functionality

var imageBuildCmd = &cobra.Command{
	Use:   "build [name]",
	Short: "Build the packer image",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		name := args[0]

		command := []string{
			"image",
			"build",
			"--template",
			fmt.Sprintf("%s.pkr.hcl", name),
		}

		runner := container.New(cfg)
		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
	simulatorCmd.AddCommand(imageCmd)
}
