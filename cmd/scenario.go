package cmd

import (
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
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
				fmt.Println("ID: " + s.Id + ", Name: " + s.DisplayName)
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
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")
			scenariosDir := viper.GetString("scenarios-dir")
			attackTag := viper.GetString("attack-container-tag")
			scenarioID := args[0]

			if err := simulator.Launch(logger, tfDir, scenariosDir, bucket, scenarioID, attackTag); err != nil {
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
