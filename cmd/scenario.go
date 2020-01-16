package cmd

import (
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	sim "github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"fmt"
)

func newScenarioListCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString("scenarios-dir")
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				logger.Errorw("Error loading scenario manifest", zap.Error(err))
				return err
			}

			fmt.Println("Available scenarios:")
			for _, s := range manifest.Scenarios {
				fmt.Println("")
				fmt.Println("Name: " + s.DisplayName)
				fmt.Println("Description: " + s.Description)
				fmt.Println("ID: " + s.Id)
			}

			return nil
		},
	}

	return cmd
}

func newScenarioLaunchCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `launch <id>`,
		Short: "Launches a scenario",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(viper.GetString("tf-dir")),
				sim.WithScenariosDir(viper.GetString("scenarios-dir")),
				sim.WithAttackTag(viper.GetString("attack-container-tag")),
				sim.WithScenarioID(args[0]),
				sim.WithBucketName(viper.GetString("state-bucket")),
				sim.WithTfVarsDir(viper.GetString("tf-vars-dir")))

			if err := simulator.Launch(); err != nil {
				if strings.HasPrefix(err.Error(), "Scenario not found") {
					logger.Warn(err.Error())
					return nil
				}
				logger.Errorw("Error launching scenario", zap.Error(err))
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

	logger, err := newLogger(viper.GetString("loglevel"), "console")
	if err != nil {
		logger.Fatalf("can't re-initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cmd.AddCommand(newScenarioListCommand(logger))
	cmd.AddCommand(newScenarioLaunchCommand(logger))

	return cmd
}
