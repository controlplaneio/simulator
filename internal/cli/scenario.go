package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/internal/container"
	"github.com/controlplaneio/simulator/scenarios"
)

var scenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Manage the simulator scenario",
}

var scenarioInstallCmd = &cobra.Command{
	Use:   "install [id]",
	Short: "Install the scenario into the simulator infrastructure",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		scenarioID := args[0]

		command := []string{
			"scenario",
			"install",
			"--name",
			scenarioID,
		}

		runner := container.New(cfg)
		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

var scenarioUninstallCmd = &cobra.Command{
	Use:   "uninstall [id]",
	Short: "Uninstall the scenario into the simulator infrastructure",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		scenarioID := args[0]

		command := []string{
			"scenario",
			"uninstall",
			"--name",
			scenarioID,
		}

		runner := container.New(cfg)
		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

var scenarioListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available scenarios",
	Run: func(cmd *cobra.Command, args []string) {
		scenarios, err := scenarios.List()
		cobra.CheckErr(err)

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

		for _, s := range scenarios {
			table.Append([]string{
				s.ID,
				s.Name,
				s.Description,
				s.Category,
				s.Difficulty})
			table.SetRowLine(true)
		}
		table.Render()
	},
}

var scenarioDescribeCmd = &cobra.Command{
	Use:   "describe [id]",
	Short: "Describes a scenario",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scenarioID := args[0]

		s, err := scenarios.Find(scenarioID)
		cobra.CheckErr(err)

		b, err := s.Challenge()
		cobra.CheckErr(err)

		fmt.Println(string(b))
	},
}

func init() {
	scenarioCmd.AddCommand(scenarioInstallCmd)
	scenarioCmd.AddCommand(scenarioUninstallCmd)
	scenarioCmd.AddCommand(scenarioListCmd)
	scenarioCmd.AddCommand(scenarioDescribeCmd)
	simulatorCmd.AddCommand(scenarioCmd)
}
