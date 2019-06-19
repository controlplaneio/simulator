package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const scenariosPath = "../simulation-scripts/scenario/"

func newScenarioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `scenario`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("scenario name is required")
			}

			name := args[0]

			scenarios, err := loadScenarios(scenariosPath)

			if err != nil {
				return err
			}

			fmt.Println("Available scenarios")
			for _, s := range scenarios {
				fmt.Println(s.DisplayName + ": " + s.Name)
			}

			if !contains(scenarios, name) {
				return fmt.Errorf("scenario %s not found", name)
			}

			fmt.Println("Chosen scenario")
			return nil
		},
	}

	return cmd
}
