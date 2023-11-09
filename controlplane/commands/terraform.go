package commands

import (
	"github.com/controlplaneio/simulator/v2/controlplane/aws"
)

type TerraformCommandType bool

const (
	Terraform        Executable           = "terraform"
	TerraformApply   TerraformCommandType = true
	TerraformDestroy TerraformCommandType = false
)

func TerraformInitCommand(workingDir string, backendConfig []string) Runnable {
	args := []string{
		"init",
		"-reconfigure",
	}

	args = append(args, backendConfig...)

	return command{
		Executable:  Terraform,
		WorkingDir:  workingDir,
		Environment: aws.Env,
		Arguments:   args,
	}
}

func TerraformCommand(workingDir string, apply TerraformCommandType, vars []string) Runnable {
	args := []string{
		"apply",
		"-auto-approve",
	}

	args = append(args, vars...)

	if !apply {
		args = append(args, "-destroy")
	}

	return command{
		Executable:  Terraform,
		WorkingDir:  workingDir,
		Environment: aws.Env,
		Arguments:   args,
	}
}
