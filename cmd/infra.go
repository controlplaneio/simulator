package cmd

import (
	sim "github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `create`,
		Short: "Runs terraform to create the required infrastructure for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucketName := viper.GetString("state-bucket")
			if bucketName == "" {
				logger.Warn("Simulator has not been initialised with an S3 bucket")
				logger.Warn("Please run simulator init")
				return nil
			}

			scenariosDir := viper.GetString("scenarios-dir")
			attackTag := viper.GetString("attack-container-tag")
			attackRepo := viper.GetString("attack-container-repo")
			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			disableIPDetection := viper.GetBool("disable-ip-detection")
			extraCIDRs := viper.GetString("extra-cidrs")

			logger.WithFields(logrus.Fields{
				"BucketName": bucketName,
			}).Info("Created s3 bucket for terraform remote state")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithScenariosDir(scenariosDir),
				sim.WithAttackTag(attackTag),
				sim.WithAttackRepo(attackRepo),
				sim.WithBucketName(bucketName),
				sim.WithoutIPDetection(disableIPDetection),
				sim.WithTfVarsDir(tfVarsDir),
				sim.WithExtraCIDRs(extraCIDRs))

			err := simulator.Create()
			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error creating infrastructure")
			}

			cfg, err := simulator.SSHConfig()
			if err != nil {
				return errors.Wrap(err, "Error getting SSH config")
			}

			err = simulator.SSHStateProvider.SaveSSHConfig(*cfg)
			if err != nil {
				return errors.Wrapf(err, "Error writing SSH config")
			}

			return err
		},
	}

	return cmd
}

func newStatusCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `status`,
		Short: "Gets the status of the infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucketName := viper.GetString("state-bucket")
			if bucketName == "" {
				logger.Warn("Simulator has not been initialised with an S3 bucket")
				logger.Warn("Please run simulator init")
				return nil
			}

			attackTag := viper.GetString("attack-container-tag")
			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			disableIPDetection := viper.GetBool("disable-ip-detection")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithAttackTag(attackTag),
				sim.WithBucketName(bucketName),
				sim.WithoutIPDetection(disableIPDetection),
				sim.WithTfVarsDir(tfVarsDir))

			tfo, err := simulator.Status()
			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error getting status of infrastructure")
				return err
			}

			if tfo.BastionPublicIP.Value == "" {
				logger.Error("No Infrastructure found")
			} else {
				logger.WithFields(logrus.Fields{
					"BastionIP": tfo.BastionPublicIP.Value,
					"MasterIP":  tfo.MasterNodesPrivateIP.Value,
					"NodeIPs":   tfo.ClusterNodesPrivateIP.Value,
				}).Info("Infrastructure Status")
			}

			return err
		},
	}

	return cmd
}

func newDestroyCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `destroy`,
		Short: "Tears down the infrastructure created for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucketName := viper.GetString("state-bucket")
			if bucketName == "" {
				logger.Warn("Simulator has not been initialised with an S3 bucket")
				logger.Warn("Please run simulator init")
				return nil
			}

			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			attackTag := viper.GetString("attack-container-tag")
			disableIPDetection := viper.GetBool("disable-ip-detection")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithAttackTag(attackTag),
				sim.WithBucketName(bucketName),
				sim.WithoutIPDetection(disableIPDetection),
				sim.WithTfVarsDir(tfVarsDir))

			err := simulator.Destroy()
			if err != nil {
				logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("Error destroying infrastructure")
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

	logger := newLogger(viper.GetString("loglevel"))

	cmd.AddCommand(newCreateCommand(logger))
	cmd.AddCommand(newStatusCommand(logger))
	cmd.AddCommand(newDestroyCommand(logger))

	return cmd
}
