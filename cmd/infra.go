package cmd

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newCreateCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `create`,
		Short: "Runs terraform to create the required infrastructure for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")
			err := simulator.Create(logger, tfDir, bucket)
			if err != nil {
				logger.Errorw("Error creating infrastructure", zap.Error(err))
			}

			return err
		},
	}

	return cmd
}

func newStatusCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `status`,
		Short: "Gets the status of the infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")
			tfo, err := simulator.Status(logger, tfDir, bucket)
			if err != nil {
				logger.Errorw("Error getting status of infrastructure", zap.Error(err))
				return err
			}

			if tfo.BastionPublicIP.Value == "" {
				logger.Error("No Infrastructure found")
			} else {
				logger.Infof("Bastion IP: %s\n", tfo.BastionPublicIP.Value)
				logger.Infof("Master IPs: %v\n", tfo.MasterNodesPrivateIP.Value)
				logger.Infof("Cluster IPs: %v\n", tfo.ClusterNodesPrivateIP.Value)
			}

			return err
		},
	}

	return cmd
}

func newDestroyCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `destroy`,
		Short: "Tears down the infrastructure created for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")

			err := simulator.Destroy(logger, tfDir, bucket)
			if err != nil {
				logger.Errorw("Error destroying infrastructure", zap.Error(err))
			}

			return err
		},
	}

	return cmd
}

func newInfraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `infra <subcommand>`,
		Short:         "Interact with AWS to create, query and destroy the required infrastructure for scenarios",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	logger, err := newLogger(viper.GetString("loglevel"), "console")
	if err != nil {
		logger.Fatalf("Can't re-initialize zap logger: %v", err)
	}

	cmd.AddCommand(newCreateCommand(logger))
	cmd.AddCommand(newStatusCommand(logger))
	cmd.AddCommand(newDestroyCommand(logger))

	return cmd
}
