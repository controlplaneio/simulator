package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/internal/config"
)

func WithConfigCmd(conf config.Config) SimulatorCmdOptions {
	var name, bucket string
	var dev, rootless bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the Simulator CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if name != "" {
				conf.Name = name
			}

			if bucket != "" {
				conf.Bucket = bucket
			}

			if dev {
				conf.Cli.Dev = true
				conf.Container.Image = "controlplane/simulator:dev"

				baseDir, err := os.Getwd()
				cobra.CheckErr(err)

				conf.BaseDir = baseDir
			} else {
				conf.Cli.Dev = false
				conf.Container.Image = "controlplane/simulator:latest"
				conf.BaseDir = ""
			}

			conf.Container.Rootless = rootless

			err := conf.Write()
			cobra.CheckErr(err)
		},
	}

	configCmd.PersistentFlags().StringVar(&name, "name", "simulator", "the name for the infrastructure")
	configCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "the s3 bucket used for storage")
	configCmd.PersistentFlags().BoolVar(&dev, "dev", false, "developer mode")
	configCmd.PersistentFlags().BoolVar(&rootless, "rootless", false, "docker running in rootless mode")

	return func(command *cobra.Command) {
		command.AddCommand(configCmd)
	}
}
