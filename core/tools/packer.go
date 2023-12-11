package tools

import (
	"context"
	"fmt"
	"io"

	"github.com/controlplaneio/simulator/v2/internal/docker"
)

const (
	PackerExecutable Executable = "packer"
)

type AMIBuilder interface {
	Build(ctx context.Context, id string) error
}

type Packer struct {
	WorkingDir string
	Output     io.Writer
}

func (p Packer) Build(ctx context.Context, id string) error {
	template := fmt.Sprintf("%s.pkr.hcl", id)

	if err := packerInitCommand(p.WorkingDir, template).Run(ctx, p.Output); err != nil {
		return fmt.Errorf("failed to initialise packer: %w", err)
	}

	if err := packerBuildCommand(p.WorkingDir, template).Run(ctx, p.Output); err != nil {
		return fmt.Errorf("failed to build ami with packer: %w", err)
	}

	return nil
}

func packerInitCommand(workingDir, template string) runner {
	args := []string{
		"init",
		template,
	}

	return runner{
		Executable: PackerExecutable,
		WorkingDir: workingDir,
		Arguments:  args,
	}
}

func packerBuildCommand(workingDir, template string) runner {
	args := []string{
		"build",
		template,
	}

	return runner{
		Executable: PackerExecutable,
		WorkingDir: workingDir,
		Arguments:  args,
	}
}

type PackerContainer struct {
	Client *docker.Client
	Config *docker.Config
}

func (p PackerContainer) Build(ctx context.Context, id string) error {
	config := *p.Config
	config.Cmd = []string{
		"ami",
		"build",
		id,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return fmt.Errorf("failed to build ami with packer: %w", err)
	}

	return nil
}
