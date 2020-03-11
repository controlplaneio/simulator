package cmd

import (
	"fmt"
	"os"
	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/kubernetes-simulator/simulator/pkg/scenario"
	sim "github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newScenarioListCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString("scenarios-dir")
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error loading scenario manifest")
				return err
			}

			difficulty := viper.GetString("difficulty")
			difficultyValues := []string{"easy", "medium", "hard", "Easy", "Medium", "Hard", ""}
			_, err = util.IsStringInSlice(difficulty, difficultyValues)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error sorting by difficulty")
				return err
			}

			category := viper.GetString("category")

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Category", "Difficulty", "Description", "ID"})
			table.SetRowLine(true)
			table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
			var output []string
			fmt.Println("Available scenarios:")
			for _, s := range manifest.Scenarios {
				if (strings.EqualFold(category, s.Category) || category == "") && (strings.EqualFold(difficulty, s.Difficulty) || difficulty == "") {
					output = []string{s.DisplayName, s.Category, s.Difficulty, s.Description, s.Id}
					table.Append(output)
				}
			}
			table.Render()
			return nil
		},
	}

	return cmd
}

func newScenarioDescribeCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `describe <id>`,
		Short: "Describes a scenario",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			scenarioID := args[0]
			manifestPath := viper.GetString("scenarios-dir")
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error loading scenario manifest")
				return err
			}

			logger.WithFields(logrus.Fields{
				"ScenarioID": scenarioID,
			}).Debug("Checking manifest contains scenario")
			if !manifest.Contains(scenarioID) {
				return errors.Errorf("Scenario not found: %s", scenarioID)
			}

			scenario := manifest.Find(scenarioID)
			err = scenario.Validate(manifestPath)
			if err != nil {
				return err
			}

			scenarioPath, err := filepath.Abs(filepath.Join(manifestPath, scenario.Path))
			if err != nil {
				return errors.Wrapf(err,
					"Error resolving %s from %s for scenario %s", scenario.Path, scenario.DisplayName, manifestPath)
			}

			challengeContent, err := ioutil.ReadFile(scenarioPath + "/challenge.txt")
			if err != nil {
				return err
			}
			challengeText := string(challengeContent)

			fmt.Println("Name: " + scenario.DisplayName + "\n")
			fmt.Println(challengeText)

			return nil
		},
	}

	return cmd
}

func newScenarioLaunchCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `launch <id>`,
		Short: "Launches a scenario",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			bucketName := viper.GetString("state-bucket")
			scenariosDir := viper.GetString("scenarios-dir")
			attackTag := viper.GetString("attack-container-tag")
			tfDir := viper.GetString("tf-dir")
			disableIPDetection := viper.GetBool("disable-ip-detection")
			tfVarsDir := viper.GetString("tf-vars-dir")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithScenariosDir(scenariosDir),
				sim.WithAttackTag(attackTag),
				sim.WithScenarioID(args[0]),
				sim.WithBucketName(bucketName),
				sim.WithoutIPDetection(disableIPDetection),
				sim.WithTfVarsDir(tfVarsDir))

			if err := simulator.Launch(); err != nil {
				if strings.HasPrefix(err.Error(), "Scenario not found") {
					logger.WithFields(logrus.Fields{
						"Error":    err,
						"Scenario": args[0],
					}).Warn("Scenario not found")

					return nil
				}
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error launching scenario")
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

	logger := newLogger(viper.GetString("loglevel"))

	cmd.AddCommand(newScenarioListCommand(logger))
	cmd.AddCommand(newScenarioLaunchCommand(logger))
	cmd.AddCommand(newScenarioDescribeCommand(logger))

	return cmd
}
