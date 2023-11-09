package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/controlplane"
)

var scenarioCmd = &cobra.Command{
	Use: "scenario",
}

var installCmd = &cobra.Command{
	Use: "install",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.InstallScenario(ctx, name)
	},
}

var uninstallCmd = &cobra.Command{
	Use: "uninstall",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		return cp.UninstallScenario(ctx, name)
	},
}

func init() {
	scenarioCmd.PersistentFlags().StringVar(&name, "name", "", "the name of the scenario to deploy")

	scenarioCmd.AddCommand(installCmd)
	scenarioCmd.AddCommand(uninstallCmd)
	simulatorCmd.AddCommand(scenarioCmd)

}
