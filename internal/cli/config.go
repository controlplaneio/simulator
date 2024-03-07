package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/internal/config"
)

func WithConfigCmd(conf config.Config) SimulatorCmdOptions {
	var name, bucket string
	var dev, rootless, printDir bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the Simulator CLI",
	}

	configCmd.PersistentFlags().StringVar(&name, "name", "simulator", "the name for the infrastructure")
	configCmd.PersistentFlags().BoolVar(&printDir, "print-dir", false, "print configuration directory")
	configCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "the s3 bucket used for storage")
	configCmd.PersistentFlags().BoolVar(&dev, "dev", false, "developer mode")
	configCmd.PersistentFlags().BoolVar(&rootless, "rootless", false, "docker running in rootless mode")

	configCmd.RunE = func(_ *cobra.Command, _ []string) error {
		if printDir {
			dir, err := config.SimulatorDir()
			if err != nil {
				return fmt.Errorf("unable to get simulator config directory: %w", err)
			}

			//nolint: forbidigo
			fmt.Println(dir)
			return nil
		}

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
			if err != nil {
				return fmt.Errorf("unable to current working directory: %w", err)
			}

			conf.BaseDir = baseDir
		} else {
			conf.Cli.Dev = false
			conf.Container.Image = "controlplane/simulator:latest"
			conf.BaseDir = ""
		}

		conf.Container.Rootless = rootless

		err := conf.Write()
		if err != nil {
			return fmt.Errorf("unable to write simulator config to disk: %w", err)
		}
		return nil
	}

	return func(command *cobra.Command) {
		command.AddCommand(configCmd)
	}
}
