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

			attackTag := viper.GetString("attack-container-tag")
			attackRepo := viper.GetString("attack-container-repo")
			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			disableIPDetection := viper.GetBool("disable-ip-detection")
			extraCIDRs := viper.GetString("extra-cidrs")
			githubUsernames := viper.GetString("github-usernames")

			logger.WithFields(logrus.Fields{
				"BucketName": bucketName,
			}).Info("Created s3 bucket for terraform remote state")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithAttackTag(attackTag),
				sim.WithAttackRepo(attackRepo),
				sim.WithBucketName(bucketName),
				sim.WithoutIPDetection(disableIPDetection),
				sim.WithTfVarsDir(tfVarsDir),
				sim.WithExtraCIDRs(extraCIDRs),
				sim.WithGithubUsernames(githubUsernames))

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

			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			disableIPDetection := viper.GetBool("disable-ip-detection")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
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
			disableIPDetection := viper.GetBool("disable-ip-detection")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
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

	cmd.PersistentFlags().StringP("tf-dir", "t", "./terraform/deployments/AWS",
		"Path to a directory containing the infrastructure scripts")
	if err := viper.BindPFlag("tf-dir", cmd.PersistentFlags().Lookup("tf-dir")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("attack-container-tag", "a", "latest",
		"The attack container tag to pull on the bastion")
	if err := viper.BindPFlag("attack-container-tag", cmd.PersistentFlags().Lookup("attack-container-tag")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("attack-container-repo", "r", "controlplane/simulator-attack",
		"The attack container repo to pull from on the bastion")
	if err := viper.BindPFlag("attack-container-repo", cmd.PersistentFlags().Lookup("attack-container-repo")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("extra-cidrs", "e", "",
		"Extra CIDRs that will be allowed to access to the bastion host. MUST be a valid CIDR and a list MUST be comma delimited")
	if err := viper.BindPFlag("extra-cidrs", cmd.PersistentFlags().Lookup("extra-cidrs")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("github-usernames", "u", "",
		"Github usernames that will be allowed access to the bastion host. MUST be a valid username and a list MUST be comma delimited")
	if err := viper.BindPFlag("github-usernames", cmd.PersistentFlags().Lookup("github-usernames")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("tf-vars-dir", "v", "/home/launch/.kubesim",
		"Path to a directory containing the terraform variables file")
	if err := viper.BindPFlag("tf-vars-dir", cmd.PersistentFlags().Lookup("tf-vars-dir")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().BoolP("disable-ip-detection", "i", false,
		"Disable public IP check. If you disable, make sure you know what you are doing.")
	if err := viper.BindPFlag("disable-ip-detection", cmd.PersistentFlags().Lookup("disable-ip-detection")); err != nil {
		panic(err)
	}

	return cmd
}
