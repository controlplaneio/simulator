package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `init`,
		Short: "Creates and configures a bucket for remote state",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			logger, err := newLogger(viper.GetString("loglevel"), "console")
			if err != nil {
				logger.Fatalf("Can't re-initialize zap logger: %v", err)
			}
			defer logger.Sync()

			if bucket == "" {
				logger.Info("No state bucket name found in config or on commandline arguments")
				logger.Debug("Asking user for a name for the state bucket")
				prompt := &survey.Input{
					Message: "Please choose a globally unique name for an S3 bucket to store the terraform state",
				}
				survey.AskOne(prompt, &bucket)
				logger.Info("Saving state bucket name to config")
				viper.Set("state-bucket", bucket)
				viper.WriteConfig()
				logger.Infof("Creating s3 bucket %s for terraform remote state\n", bucket)
				simulator.CreateRemoteStateBucket(logger, bucket)
				logger.Infof("Created s3 bucket %s for terraform remote state\n", bucket)
				return nil
			}

			logger.Warnf("Simulator is already configured to use an S3 bucket named %s", bucket)
			logger.Warn("Please remove the state-bucket from simulator.yaml to create another")
			return nil
		},
	}

	return cmd
}
