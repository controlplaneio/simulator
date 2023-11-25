package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/tools"
	"github.com/controlplaneio/simulator/v2/scenarios"
)

func WithScenarioCmd(opts ...SimulatorCmdOptions) SimulatorCmdOptions {
	amiCmd := &cobra.Command{
		Use:   "scenario",
		Short: "Manage the simulator scenarios",
	}

	for _, opt := range opts {
		opt(amiCmd)
	}

	return func(command *cobra.Command) {
		command.AddCommand(amiCmd)
	}
}

func WithScenarioListCmd() SimulatorCmdOptions {
	scenarioListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available scenarios",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := scenarios.List()
			cobra.CheckErr(err)

			tabulateScenarios(list)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(scenarioListCmd)
	}
}

func WithScenarioDescribeCmd() SimulatorCmdOptions {
	scenarioDescribeCmd := &cobra.Command{
		Use:   "describe [id]",
		Short: "Describes a scenario",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			scenarioID := args[0]

			s, err := scenarios.Find(scenarioID)
			cobra.CheckErr(err)

			tabulateScenarios([]scenarios.Scenario{s})
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(scenarioDescribeCmd)
	}
}

func WithScenarioInstallCmd(manager tools.ScenarioManager) SimulatorCmdOptions {
	scenarioInstallCmd := &cobra.Command{
		Use:   "install [id]",
		Short: "Install the scenario into the simulator infrastructure",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			scenarioID := args[0]
			err := manager.Install(ctx, scenarioID)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(scenarioInstallCmd)
	}
}

func WithScenarioUninstallCmd(manager tools.ScenarioManager) SimulatorCmdOptions {
	scenarioInstallCmd := &cobra.Command{
		Use:   "uninstall [id]",
		Short: "Uninstall the scenario from the simulator infrastructure",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			scenarioID := args[0]
			err := manager.Uninstall(ctx, scenarioID)
			cobra.CheckErr(err)
		},
	}

	return func(command *cobra.Command) {
		command.AddCommand(scenarioInstallCmd)
	}
}

func tabulateScenarios(scenarios []scenarios.Scenario) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{
		"ID",
		"Name",
		"Description",
		"Category",
		"Difficulty",
	})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	for _, scenario := range scenarios {
		table.Append([]string{
			scenario.ID,
			scenario.Name,
			scenario.Description,
			scenario.Category,
			scenario.Difficulty,
		})
		table.SetRowLine(true)
	}
	table.Render()
}
