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

	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newInfraCommand())
	cmd.AddCommand(newScenarioCommand())
	cmd.AddCommand(newSSHCommand())
	cmd.AddCommand(newVersionCommand())

	cmd.PersistentFlags().StringP("bucket", "b", "",
		"The name of the s3 bucket to use.  Must be globally unique and ill be prefixed with 'simulator-'")
	cmd.MarkFlagRequired("bucket")
	viper.BindPFlag("bucket", cmd.PersistentFlags().Lookup("bucket"))

	cmd.PersistentFlags().StringP("loglevel", "l", "info", "Level of detail in output logging")
	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))

	cmd.PersistentFlags().StringP("tf-dir", "t", "./terraform/deployments/AWS",
		"Path to a directory containing the infrastructure scripts")
	viper.BindPFlag("tf-dir", cmd.PersistentFlags().Lookup("tf-dir"))

	// TODO: (rem) this is also used to locate the perturb.sh script which may be subsumed by this app
	cmd.PersistentFlags().StringP("scenarios-dir", "s", "./simulation-scripts",
		"Path to a directory containing a scenario manifest")
	viper.BindPFlag("scenarios-dir", cmd.PersistentFlags().Lookup("scenarios-dir"))

	return cmd
}

func initConfig() {
	viper.SetConfigType("yaml")
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
