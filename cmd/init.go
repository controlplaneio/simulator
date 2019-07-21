package cmd

import (
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

			logger.Infof("Creating s3 bucket %s for terraform remote state\n", bucket)

			return nil
		},
	}

	return cmd
}
