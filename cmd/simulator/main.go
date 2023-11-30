package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/aws"
	"github.com/controlplaneio/simulator/v2/core/tools"
	"github.com/controlplaneio/simulator/v2/internal/cli"
	"github.com/controlplaneio/simulator/v2/internal/config"
	"github.com/controlplaneio/simulator/v2/internal/docker"
	"github.com/controlplaneio/simulator/v2/internal/logging"
)

const (
	ownerReadWriteExecute = 0700
)

func main() {
	logging.Configure()

	conf := config.Config{}
	if err := conf.Read(); err != nil {
		slog.Error("failed to read config", "error", err)
		os.Exit(1)
	}

	adminBundleDir, err := conf.AdminBundleDir()
	if err != nil {
		slog.Error("failed to determine admin bundle dir", "error", err)
		os.Exit(1)
	}

	playerBundleDir, err := conf.PlayerBundleDir()
	if err != nil {
		slog.Error("failed to determine player bundle dir", "error", err)
		os.Exit(1)
	}

	mkDirsIfNotExist(adminBundleDir, playerBundleDir)

	mounts := []docker.MountConfig{
		{
			Source: adminBundleDir,
			Target: "/simulator/config/admin",
		},
		{
			Source: playerBundleDir,
			Target: "/simulator/config/player",
		},
		{
			Source:   "/home/ric/.aws",
			Target:   "/home/ubuntu/.aws",
			ReadOnly: true,
		},
	}

	// If running in dev mode, mount the configuration directories
	if conf.Cli.Dev {
		mounts = append(mounts, []docker.MountConfig{
			{
				Source: filepath.Join(conf.BaseDir, "packer"),
				Target: "/simulator/packer",
			},
			{
				Source: filepath.Join(conf.BaseDir, "terraform"),
				Target: "/simulator/terraform",
			},
			{
				Source: filepath.Join(conf.BaseDir, "ansible"),
				Target: "/simulator/ansible",
			},
		}...)
	}

	dockerConfig := &docker.Config{
		Image:    conf.Container.Image,
		Rootless: conf.Container.Rootless,
		Env:      aws.EnvVars(),
		Mounts:   mounts,
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		slog.Error("failed to create docker client", "error", err)
		os.Exit(1)
	}

	awsBucketCreator, err := aws.NewS3Client(context.Background())
	if err != nil {
		slog.Error("failed to create s3 client", "error", err)
		os.Exit(1)
	}

	amiManager := aws.EC2{}
	amiCreator := tools.PackerContainer{
		Client: dockerClient,
		Config: dockerConfig,
	}
	infraManager := tools.TerraformContainer{
		Client: dockerClient,
		Config: dockerConfig,
	}
	scenarioManager := tools.AnsiblePlaybookContainer{
		Client: dockerClient,
		Config: dockerConfig,
	}

	withStateBucketFlag := cli.WithFlag("stateBucket", conf.Bucket, "the name of the S3 bucket to store Terraform state")
	withStateKeyFlag := cli.WithFlag("stateKey", "terraform.tfstate", "the path to the state file in the S3 bucket")
	withNameFlag := cli.WithFlag("name", conf.Name, "the name used for the Simulator infrastructure")

	simulator := cli.NewSimulatorCmd(
		cli.WithConfigCmd(conf),
		cli.WithBucketCmd(
			cli.WithCreateBucketCmd(conf, awsBucketCreator),
		),
		cli.WithContainerCmd(
			cli.WithContainerPullCmd(conf, dockerClient),
		),
		cli.WithAMICmd(
			cli.WithAmiBuildCmd(amiCreator),
			cli.WithAMIListCmd(amiManager),
			cli.WithAMIDeleteCmd(amiManager),
		),
		cli.WithInfraCmd(
			cli.WithInfraCreateCmd(infraManager,
				withStateBucketFlag,
				withStateKeyFlag,
				withNameFlag,
			),
			cli.WithInfraDestroyCmd(infraManager,
				withStateBucketFlag,
				withStateKeyFlag,
				withNameFlag,
			),
		),
		cli.WithScenarioCmd(
			cli.WithScenarioListCmd(),
			cli.WithScenarioDescribeCmd(),
			cli.WithScenarioInstallCmd(scenarioManager),
			cli.WithScenarioUninstallCmd(scenarioManager),
		),
	)

	err = simulator.Execute()
	cobra.CheckErr(err)
}

func mkDirsIfNotExist(dirs ...string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(dir, ownerReadWriteExecute); err != nil {
				slog.Error("failed to bundle directory", "dir", dir, "error", err)
				os.Exit(1)
			}
		}
	}
}
