package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/core/tools"
	"github.com/controlplaneio/simulator/v2/internal/cli"
	"github.com/controlplaneio/simulator/v2/internal/config"
)

func main() {
	// create variables for the file paths within the container
	simulatorDir := "/simulator"
	adminBundleDir := filepath.Join(simulatorDir, "config", "admin")
	packerDir := filepath.Join(simulatorDir, "packer")
	terraformWorkspaceDir := filepath.Join(simulatorDir, "terraform/workspaces/simulator")
	ansibleConfigPath := filepath.Join(adminBundleDir, "ansible.cfg")
	ansiblePlaybookDir := filepath.Join(simulatorDir, "ansible/playbooks")

	conf := config.Config{}
	if err := conf.Read(); err != nil {
		slog.Error("failed to read config", "error", err)
		os.Exit(1)
	}

	amiBuilder := tools.Packer{
		WorkingDir: packerDir,
		Output:     os.Stdout,
	}

	infraManager := tools.Terraform{
		WorkingDir: terraformWorkspaceDir,
		Output:     os.Stdout,
	}

	scenarioManager := tools.AnsiblePlaybook{
		WorkingDir:  adminBundleDir,
		PlaybookDir: ansiblePlaybookDir,
		// Ansible complains on Windows+WSL that the directory ansible configuration is world writable
		// and hence ignore the configuration unless explicitly set using the ANSIBLE_CONFIG environment variable.
		Env:    []string{"ANSIBLE_CONFIG=" + ansibleConfigPath},
		Output: os.Stdout,
	}

	withStateBucketFlag := cli.WithFlag("stateBucket", "", "the name of the S3 bucket to store Terraform state")
	withStateKeyFlag := cli.WithFlag("stateKey", "", "the path to the state file in the S3 bucket")
	withNameFlag := cli.WithFlag("name", "", "the name used for the Simulator infrastructure")

	simulator := cli.NewSimulatorCmd(
		cli.WithAMICmd(
			cli.WithAmiBuildCmd(amiBuilder),
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
			cli.WithScenarioInstallCmd(scenarioManager),
			cli.WithScenarioUninstallCmd(scenarioManager),
		),
	)

	err := simulator.Execute()
	cobra.CheckErr(err)
}
