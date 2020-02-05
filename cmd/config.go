package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigGetCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `get <key>`,
		Short: "Gets the value of a setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				logger.Fatal("<key> is required")
				return nil
			}

			key := args[0]
			logger.Infof(`%s = %v\n`, key, viper.Get(key))

			return nil
		},
	}

	return cmd
}

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `config`,
		Short:         "Interact with simulator config",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	logger := newLogger(viper.GetString("loglevel"))
	cmd.AddCommand(newConfigGetCommand(logger))

	return cmd
}
