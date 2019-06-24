package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator",
		Short: "Simulator command line",
		Long: `
A distributed systems and infrastructure simulator for attacking and
debugging Kubernetes
`,
	}

	cmd.AddCommand(newInfraCommand())
	cmd.AddCommand(newScenarioCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

func Execute() error {
	cmd := NewCmdRoot()
	return cmd.Execute()
}
