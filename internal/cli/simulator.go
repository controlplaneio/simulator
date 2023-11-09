package cli

import (
	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/internal/config"
)

var cfg *config.Config

var simulatorCmd = &cobra.Command{
	Use:   "simulator",
	Short: "Simulator CLI",
}

func Execute() {
	err := simulatorCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	cfg = &config.Config{}

	cobra.OnInitialize(readConfig)
	cobra.OnFinalize(writeConfig)
}

func readConfig() {
	err := cfg.Read
	cobra.CheckErr(err())
}

func writeConfig() {
	err := cfg.Write
	cobra.CheckErr(err())
}
