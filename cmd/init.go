package cmd

import (
	"bufio"
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strings"
)

func saveBucketConfig(logger *zap.SugaredLogger, bucket string) {
	logger.Info("Saving state bucket name to config")
	viper.Set("state-bucket", bucket)
	viper.WriteConfig()
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `init`,
		Short: "Creates and configures a bucket for remote state",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")
			tfDir := viper.GetString("tf-dir")

			logger, err := newLogger(viper.GetString("loglevel"), "console")
			if err != nil {
				logger.Fatalf("Can't re-initialize zap logger: %v", err)
			}
			defer logger.Sync()

			if bucket == "" {
				logger.Info("No state bucket name found in config or on commandline arguments")
				logger.Debug("Asking user for a name for the state bucket")

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Please choose a globally unique name for an S3 bucket to store the terraform state: ")
				bucket, err := reader.ReadString('\n')
				if err != nil {
					return errors.Wrap(err, "Error reading bucket nbame from stdin")
				}

				bucket = strings.TrimSpace(bucket)


				logger.Infof("Creating s3 bucket %s for terraform remote state\n", bucket)
				if err = simulator.CreateRemoteStateBucket(logger, bucket); err != nil {
					if strings.HasPrefix(errors.Cause(err).Error(), "BucketAlreadyOwnedByYou") {
						logger.Infof("%s already exists and you own it", bucket)
						saveBucketConfig(logger, bucket)
						return nil
					}

					return errors.Wrapf(err, "Error creating s3 bucket %s", bucket)
				}
				saveBucketConfig(logger, bucket)

				return nil
			}

			logger.Warnf("Simulator is already configured to use an S3 bucket named %s", bucket)
			logger.Warn("Please remove the state-bucket from simulator.yaml to create another")
			return nil
		},
	}

	return cmd
}
