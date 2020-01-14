package cmd

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newSSHConfigCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `config`,
		Short: "Prints the stanzas to add to ssh config to connect to your cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			simulator := simulator.NewSimulator(
				simulator.WithLogger(logger),
				simulator.WithTfDir(viper.GetString("tf-dir")),
				simulator.WithScenariosDir(viper.GetString("scenarios-dir")),
				simulator.WithAttackTag(viper.GetString("attack-container-tag")),
				simulator.WithBucketName(viper.GetString("state-bucket")),
				simulator.WithTfVarsDir(viper.GetString("tf-vars-dir")))

			cfg, err := simulator.SSHConfig()
			if err != nil {
				return errors.Wrap(err, "Error getting SSH config")
			}

			shouldwrite, err := cmd.Flags().GetBool("write")
			if err != nil {
				return errors.Wrapf(err, "Error getting --write cli flag")
			}

			if !shouldwrite {
				logger.Info(*cfg)
			}

			err = ssh.EnsureSSHConfig(*cfg)
			if err != nil {
				return errors.Wrapf(err, "Error writing SSH config")
			}

			return nil
		},
	}
	cmd.PersistentFlags().Bool("write", true, "Write the ssh config - this is isolated from your default SSH config")

	return cmd
}

func newSSHAttackCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `attack`,
		Short: "Connect to an attack container to complete the scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")
			tfVarsDir := viper.GetString("tf-vars-dir")
			attackTag := viper.GetString("attack-container-tag")

			return simulator.Attack(logger, tfDir, bucket, attackTag, tfVarsDir)
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

	logger, err := newLogger(viper.GetString("loglevel"), "console")
	if err != nil {
		logger.Fatalf("can't re-initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cmd.AddCommand(newSSHConfigCommand(logger))
	cmd.AddCommand(newSSHAttackCommand(logger))

	return cmd
}
