package cmd

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/runner"
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newScenarioListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := scenario.ManifestPath()
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				return err
			}

			fmt.Println("Available scenarios:")
			for _, s := range manifest.Scenarios {
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

			scenarioID := args[0]

			manifestPath := scenario.ManifestPath()
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				return err
			}

			if !manifest.Contains(scenarioID) {
				return fmt.Errorf("scenario %s not found", scenarioID)
			}

			tfo, err := runner.Status()
			if !tfo.IsUsable() {
				return fmt.Errorf("No infrastructure, please run simulator infra create:\n %#v", tfo)
			}

			scenarioPath := manifest.Find(scenarioID).Path

			po := runner.MakePerturbOptions(*tfo, scenarioPath)
			fmt.Println("Converted usable terraform output into perturb options")
			fmt.Printf("%#v", po)
			c, err := tfo.ToSSHConfig()
			if err != nil {
				return err
			}

			cp, err := util.ExpandTilde("~/.ssh/config")
			if err != nil {
				return err
			}

			written, err := util.EnsureFile(*cp, *c)
			if err != nil {
				return err
			}

			if !written {
				fmt.Printf("Please add the following lines to your ssh config\n---\n%s\n---\n", *c)
			}

			_, err = runner.Perturb(&po)
			if err != nil {
				return errors.Wrapf(err, "Error running perturb with %#v", po)
			}

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
