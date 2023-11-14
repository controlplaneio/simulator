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
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.InstallScenario(ctx, name)
		cobra.CheckErr(err)
	},
}

var uninstallCmd = &cobra.Command{
	Use: "uninstall",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		cp := controlplane.New()
		err := cp.UninstallScenario(ctx, name)
		cobra.CheckErr(err)
	},
}

func init() {
	scenarioCmd.PersistentFlags().StringVar(&name, "name", "", "the name of the scenario to deploy")

	scenarioCmd.AddCommand(installCmd)
	scenarioCmd.AddCommand(uninstallCmd)
	simulatorCmd.AddCommand(scenarioCmd)
}
