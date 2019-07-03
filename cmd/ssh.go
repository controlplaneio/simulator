package cmd

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
)

func newSSHConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `config`,
		Short: "Prints the stanzas to add to ssh config to connect to your cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := scenario.ManifestPath()
			cfg, err := simulator.Config(p)
			if err != nil {
				return err
			}

			fmt.Println(*cfg)

			return nil
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

	cmd.AddCommand(newSSHConfigCommand())

	return cmd
}
