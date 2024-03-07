package cli

import (
	"context"
	"fmt"
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
		RunE: func(_ *cobra.Command, _ []string) error {
			list, err := scenarios.List()
			if err != nil {
				return fmt.Errorf("unable to list available scenarios: %w", err)
			}

			tabulateScenarios(list)
			return nil
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
		RunE: func(_ *cobra.Command, args []string) error {
			scenarioID := args[0]

			s, err := scenarios.Find(scenarioID)
			if err != nil {
				return fmt.Errorf("unable to describe scenario: %w", err)
			}

			tabulateScenarios([]scenarios.Scenario{s})
			return nil
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
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			scenarioID := args[0]
			err := manager.Install(ctx, scenarioID)
			if err != nil {
				return fmt.Errorf("unable to install the scenario: %w", err)
			}
			return nil
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
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			scenarioID := args[0]
			err := manager.Uninstall(ctx, scenarioID)
			if err != nil {
				return fmt.Errorf("unable to uninstall the scenario: %w", err)
			}
			return nil
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
