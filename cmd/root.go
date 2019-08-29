package cmd

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "simulator",
	Short: "Simulator command line",
	Long: `
A distributed systems and infrastructure simulator for attacking and
debugging Kubernetes
`,
}

var logger *zap.SugaredLogger

func newCmdRoot() *cobra.Command {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "c", "", "Path to the simulator config file")
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(newConfigCommand())
	rootCmd.AddCommand(newInfraCommand())
	rootCmd.AddCommand(newInitCommand())
	rootCmd.AddCommand(newScenarioCommand())
	rootCmd.AddCommand(newSSHCommand())
	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(newCompletionCmd())

	rootCmd.PersistentFlags().StringP("state-bucket", "b", "",
		"The name of the s3 bucket to use for remote-state.  Must be globally unique")
	rootCmd.MarkFlagRequired("state-bucket")
	viper.BindPFlag("state-bucket", rootCmd.PersistentFlags().Lookup("state-bucket"))

	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "Level of detail in output logging")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().StringP("tf-dir", "t", "./terraform/deployments/AWS",
		"Path to a directory containing the infrastructure scripts")
	viper.BindPFlag("tf-dir", rootCmd.PersistentFlags().Lookup("tf-dir"))

	// TODO: (rem) this is also used to locate the perturb.sh script which may be
	// subsumed by this app
	rootCmd.PersistentFlags().StringP("scenarios-dir", "s", "./simulation-scripts",
		"Path to a directory containing a scenario manifest")
	viper.BindPFlag("scenarios-dir", rootCmd.PersistentFlags().Lookup("scenarios-dir"))

	return rootCmd
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
		// todo(ajm) this errors if not in the same dir as `simulator.yaml`. Move
		// those vars here?
		panic(errors.Wrapf(err, "Error reading config file"))
	}

	// read config from environment too
	viper.SetEnvPrefix("simulator")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates Bash completion scripts",
		Long: `To load completion run

. <(simulator completion)

To configure your Bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(simulator completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
	return cmd
}

// Execute starts the aplication
func Execute() error {
	cmd := newCmdRoot()

	return cmd.Execute()
}
