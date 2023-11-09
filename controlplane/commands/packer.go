package commands

import (
	"github.com/controlplaneio/simulator/v2/controlplane/aws"
)

const (
	Packer Executable = "packer"
)

func PackerInitCommand(workingDir, template string) Runnable {
	args := []string{
		"init",
		template,
	}

	return command{
		Executable:  Packer,
		WorkingDir:  workingDir,
		Environment: aws.Env,
		Arguments:   args,
	}
}

func PackerBuildCommand(workingDir, template string) Runnable {
	args := []string{
		"build",
		template,
	}

	return command{
		Executable:  Packer,
		WorkingDir:  workingDir,
		Environment: aws.Env,
		Arguments:   args,
	}
}
