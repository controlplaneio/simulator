package cmd

import (
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newScenarioListCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString("scenarios-dir")
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				logger.Errorw("Error loading scenario manifest", err)
				return err
			}

			logger.Info("Available scenarios:")
			for _, s := range manifest.Scenarios {
				logger.Info("ID: " + s.Id + ", Name: " + s.DisplayName)
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
			if len(args) < 1 {
				logger.Fatalf("scenario id is required")
				return nil
			}

			bucket := viper.GetString("bucket")
			tfDir := viper.GetString("tf-dir")
			scenariosDir := viper.GetString("scenarios-dir")
			scenarioID := args[0]

			err := simulator.Launch(logger, tfDir, scenariosDir, bucket, scenarioID)
			if err != nil {
				logger.Errorw("Error launching scenario", err)
			}

			return err
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
