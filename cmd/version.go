package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// vars injected by goreleaser at build time
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `version`,
		Short: "Prints simulator version",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := newLogger(viper.GetString("loglevel"), "console")
			if err != nil {
				logger.Fatalf("Can't re-initialize zap logger: %v", err)
			}
			defer logger.Sync()

			logger.Infof("version %s\n", version)
			logger.Infof("git commit %s\n", commit)
			logger.Infof("build date %s\n", date)

			return nil
		},
	}

	return cmd
}
