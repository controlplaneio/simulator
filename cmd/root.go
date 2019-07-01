package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator",
		Short: "Simulator command line",
		Long: `
A distributed systems and infrastructure simulator for attacking and
debugging Kubernetes
`,
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "c", "", "the directory where simulator.yaml can be found")
	cobra.OnInitialize(initConfig)

	cmd.AddCommand(newInfraCommand())
	cmd.AddCommand(newScenarioCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newVersionCommand())

	cmd.PersistentFlags().StringP("loglevel", "l", "info", "the level of detail in output logging")
	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))

	return cmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("simulator")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.Wrapf(err, "Error reading config file at ./simulator.yaml"))
	}

	// read config from environment too
	viper.SetEnvPrefix("simulator")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

func Execute() error {
	cmd := NewCmdRoot()
	return cmd.Execute()
}
