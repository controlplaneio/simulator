package cmd

import (
	"bufio"
	"fmt"
	"github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func saveBucketConfig(logger *logrus.Logger, bucket string) {
	logger.Info("Saving state bucket name to config")
	viper.Set("state-bucket", bucket)
	if err := viper.WriteConfig(); err != nil {
		logger.WithFields(logrus.Fields{
			"Error":      err,
			"BucketName": bucket,
		}).Fatal("Unable to write config")
	}
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `init`,
		Short: "Creates and configures a bucket for remote state",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("state-bucket")

			logger := newLogger(viper.GetString("loglevel"))

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

				logger.WithFields(logrus.Fields{
					"BucketName": bucket,
				}).Info("Creating s3 bucket for terraform remote state")
				if err = simulator.CreateRemoteStateBucket(logger, bucket); err != nil {
					if strings.HasPrefix(errors.Cause(err).Error(), "BucketAlreadyOwnedByYou") {
						logger.WithFields(logrus.Fields{
							"BucketName": bucket,
						}).Info("Bucket already exists and you own it")
						saveBucketConfig(logger, bucket)
						return nil
					}

					return errors.Wrapf(err, "Error creating s3 bucket %s", bucket)
				}
				saveBucketConfig(logger, bucket)

				return nil
			}

			logger.WithFields(logrus.Fields{
				"BucketName": bucket,
			}).Warn("Simulator is already configured to use an S3 bucket")
			logger.Warn("Please remove the state-bucket from simulator.yaml to create another")
			return nil
		},
	}

	return cmd
}
