package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// vars injected by goreleaser at build time
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func newVersionCommand(logger *zap.SugaredLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `version`,
		Short: "Prints simulator version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("version %s\ngit commit %s\nbuild date %s\n", version, commit, date)
			return nil
		},
	}

	return cmd
}
