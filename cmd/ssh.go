package cmd

import (
	sim "github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSSHConfigCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `config`,
		Short: "Prints the stanzas to add to ssh config to connect to your cluster",
		RunE: func(cmd *cobra.Command, args []string) error {

			bucketName := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")

			simulator := sim.NewSimulator(
				sim.WithLogger(logger),
				sim.WithTfDir(tfDir),
				sim.WithBucketName(bucketName),
				sim.WithTfVarsDir(tfVarsDir))

			cfg, err := simulator.SSHConfig()
			if err != nil {
				return errors.Wrap(err, "Error getting SSH config")
			}

			shouldwrite, err := cmd.Flags().GetBool("write")
			if err != nil {
				return errors.Wrapf(err, "Error getting --write cli flag")
			}

			if !shouldwrite {
				logger.Info(cfg)
			}

			err = simulator.SSHStateProvider.SaveSSHConfig(*cfg)
			if err != nil {
				return errors.Wrapf(err, "Error writing SSH config")
			}

			return nil
		},
	}
	cmd.PersistentFlags().Bool("write", true, "Write the ssh config - this is isolated from your default SSH config")

	return cmd
}

func newSSHAttackCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `attack`,
		Short: "Connect to an attack container to complete the scenario",
		RunE: func(cmd *cobra.Command, args []string) error {

			bucketName := viper.GetString("state-bucket")
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

			return simulator.Attack()
		},
	}

	return cmd

}

func newSSHCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `ssh <command>`,
		Short:         "Interact with the cluster",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.PersistentFlags().StringP("tf-vars-dir", "v", "/home/launch/.kubesim",
		"Path to a directory containing the terraform variables file")
	if err := viper.BindPFlag("tf-vars-dir", rootCmd.PersistentFlags().Lookup("tf-vars-dir")); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("tf-dir", "t", "./terraform/deployments/AWS",
		"Path to a directory containing the infrastructure scripts")
	if err := viper.BindPFlag("tf-dir", rootCmd.PersistentFlags().Lookup("tf-dir")); err != nil {
		panic(err)
	}

	logger := newLogger(viper.GetString("loglevel"))

	cmd.AddCommand(newSSHConfigCommand(logger))
	cmd.AddCommand(newSSHAttackCommand(logger))

	return cmd
}
