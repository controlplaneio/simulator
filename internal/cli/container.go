package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/internal/config"
	"github.com/controlplaneio/simulator/v2/internal/docker"
)

func WithContainerCmd(opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	containerCmd := &cobra.Command{
		Use:   "container",
		Short: "Manage Simulator Container images",
	}

	for _, opt := range opts {
		opt(containerCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(containerCmd)
	}
}

func WithContainerPullCmd(config config.Config, client *docker.Client) SimulatorCmdOptions {
	imagePullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull the Simulator Container Image",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := client.PullImage(ctx, config.Container.Image)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(imagePullCmd)
	}
}
