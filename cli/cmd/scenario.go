package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newScenarioListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := manifestPath()
			scenarios, err := loadScenarios(manifestPath)

			if err != nil {
				return err
			}

			fmt.Println("Available scenarios:")
			for _, s := range scenarios {
				fmt.Println("ID: " + s.Id + ", Name: " + s.DisplayName)
			}

			return nil
		},
	}

	return cmd
}

func newScenarioLaunchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `launch <id>`,
		Short: "Launches a scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("scenario id is required")
			}

			scenarioId := args[0]

			manifestPath := manifestPath()
			scenarios, err := loadScenarios(manifestPath)

			if err != nil {
				return err
			}

			if !contains(scenarios, scenarioId) {
				return fmt.Errorf("scenario %s not found", scenarioId)
			}

			fmt.Println("Chosen scenario: " + scenarioId)
			return nil
		},
	}

	return cmd
}

func newScenarioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `scenario <subcommand>`,
		Short:         "Interact with scenarios",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(newScenarioListCommand())
	cmd.AddCommand(newScenarioLaunchCommand())

	return cmd
}
