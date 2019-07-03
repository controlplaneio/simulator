package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `get <key>`,
		Short: "Gets the value of a setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("<key> is required")
			}

			key := args[0]
			fmt.Printf(`%s = %v\n`, key, viper.Get(key))

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

	cmd.AddCommand(newConfigGetCommand())

	return cmd
}
