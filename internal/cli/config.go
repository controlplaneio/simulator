package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var name, bucket string
var dev bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the simulator cli",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, err := os.Getwd()
		cobra.CheckErr(err)

		cfg.BaseDir = baseDir

		if name != "" {
			cfg.Name = name
		}

		if bucket != "" {
			cfg.Bucket = bucket
		}

		if dev {
			cfg.Cli.Dev = true
			cfg.Container.Image = "controlplane/simulator:dev"
		} else {
			cfg.Cli.Dev = false
			cfg.Container.Image = "controlplane/simulator:latest"
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringVar(&name, "name", "simulator", "the name for the infrastructure")
	configCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "the s3 bucket used for storage")
	configCmd.PersistentFlags().BoolVar(&dev, "dev", false, "developer mode")

	simulatorCmd.AddCommand(configCmd)
}
