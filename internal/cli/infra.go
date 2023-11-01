package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/internal/container"
)

var infraCmd = &cobra.Command{
	Use:   "infra [command]",
	Short: "Manage the simulator infrastructure",
}

var infraCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create simulator infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		runner := container.New(cfg)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		command := []string{
			"infra",
			"create",
			"--bucket",
			cfg.Bucket,
			"--name",
			cfg.Name,
		}

		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

var infraDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy simulator infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		runner := container.New(cfg)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		command := []string{
			"infra",
			"destroy",
			"--bucket",
			cfg.Bucket,
			"--name",
			cfg.Name,
		}

		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

func init() {
	infraCmd.AddCommand(infraCreateCmd)
	infraCmd.AddCommand(infraDestroyCmd)

	simulatorCmd.AddCommand(infraCmd)
}
