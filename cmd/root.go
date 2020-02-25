package cmd

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var logger *logrus.Logger //nolint:deadcode,unused

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

	// NOTE the panics here are needed - if these calls fail we cannot recover
	// and the cause is most likely programmer error
	rootCmd.PersistentFlags().StringP("state-bucket", "b", "",
		"The name of the s3 bucket to use for remote-state.  Must be globally unique")
	if err := viper.BindPFlag("state-bucket", rootCmd.PersistentFlags().Lookup("state-bucket")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("difficulty", "d", "",
		"Sorts the list of scenarios by only showing scenarios of specified difficulty")
	if err := viper.BindPFlag("difficulty", rootCmd.PersistentFlags().Lookup("difficulty")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("category", "g", "",
		"Sorts the list of scenarios by only showing scenarios of specified category")
	if err := viper.BindPFlag("category", rootCmd.PersistentFlags().Lookup("category")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "Level of detail in output logging")
	if err := viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("tf-dir", "t", "./terraform/deployments/AWS",
		"Path to a directory containing the infrastructure scripts")
	if err := viper.BindPFlag("tf-dir", rootCmd.PersistentFlags().Lookup("tf-dir")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("attack-container-tag", "a", "latest",
		"The attack container tag to pull on the bastion")
	if err := viper.BindPFlag("attack-container-tag", rootCmd.PersistentFlags().Lookup("attack-container-tag")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("attack-container-repo", "r", "controlplane/simulator-attack",
		"The attack container repo to pull from on the bastion")
	if err := viper.BindPFlag("attack-container-repo", rootCmd.PersistentFlags().Lookup("attack-container-repo")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("extra-cidrs", "e", "",
		"Extra CIDRs that will be allowed to access to the bastion host. MUST be a valid CIDR and a list MUST be comma delimited")
	if err := viper.BindPFlag("extra-cidrs", rootCmd.PersistentFlags().Lookup("extra-cidrs")); err != nil {
		panic(err)
	}

	// TODO: (rem) this is also used to locate the perturb.sh script which may be
	// subsumed by this app
	rootCmd.PersistentFlags().StringP("scenarios-dir", "s", "./simulation-scripts",
		"Path to a directory containing a scenario manifest")
	if err := viper.BindPFlag("scenarios-dir", rootCmd.PersistentFlags().Lookup("scenarios-dir")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("tf-vars-dir", "v", "/home/launch/.kubesim",
		"Path to a directory containing the terraform variables file")
	if err := viper.BindPFlag("tf-vars-dir", rootCmd.PersistentFlags().Lookup("tf-vars-dir")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().BoolP("disable-ip-detection", "i", false,
		"Disable public IP check. If you disable, make sure you know what you are doing.")
	if err := viper.BindPFlag("disable-ip-detection", rootCmd.PersistentFlags().Lookup("disable-ip-detection")); err != nil {
		panic(err)
	}

	return rootCmd
}

func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.kubesim")
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
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				panic(err)
			}
		},
	}
	return cmd
}

// Execute starts the aplication
func Execute() error {
	cmd := newCmdRoot()

	return cmd.Execute()
}
