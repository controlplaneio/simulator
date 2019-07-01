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

	cmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "c", "", "Path to the simulator config file")
	cobra.OnInitialize(initConfig)

	cmd.AddCommand(newInfraCommand())
	cmd.AddCommand(newScenarioCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newVersionCommand())

	cmd.PersistentFlags().StringP("loglevel", "l", "info", "Level of detail in output logging")
	cmd.PersistentFlags().StringP("tf-dir", "t", "./terraform", "Path to a directory containing the infrastructure scripts")
	// TODO: (rem) this is also used to locate the perturb.sh script
	cmd.PersistentFlags().StringP("scenarios-dir", "s", "./simulation-scripts", "Path to a directory containing a scenario manifest")
	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("tf-dir", cmd.PersistentFlags().Lookup("tf-dir"))
	viper.BindPFlag("scenarios-dir", cmd.PersistentFlags().Lookup("scenarios-dir"))

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
		panic(errors.Wrapf(err, "Error reading config file"))
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
