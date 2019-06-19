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
			scenarios, err := loadScenarios(scenariosPath)

			if err != nil {
				return err
			}

			for _, s := range scenarios {
				fmt.Println(s.DisplayName)
			}

			return nil
		},
	}

	return cmd
}
