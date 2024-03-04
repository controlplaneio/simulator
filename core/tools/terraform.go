package tools

import (
	"context"
	"fmt"
	"io"

	"github.com/controlplaneio/simulator/v2/internal/docker"
)

type TerraformCommandType bool

const (
	TerraformExecutable Executable           = "terraform"
	TerraformApply      TerraformCommandType = true
	TerraformDestroy    TerraformCommandType = false
)

type InfraManager interface {
	Create(ctx context.Context, stateBucket string, stateKey string, name string) error
	Destroy(ctx context.Context, stateBucket string, stateKey string, name string) error
}

type Terraform struct {
	WorkingDir string
	StdOut     io.Writer
	StdErr     io.Writer
}

func (t Terraform) Create(ctx context.Context, stateBucket string, stateKey string, name string) error {
	backend := backendConfig(stateBucket, stateKey)

	err := terraformInitCommand(t.WorkingDir, backend).Run(ctx, t.StdOut, t.StdErr)
	if err != nil {
		return fmt.Errorf("failed to initialise terraform: %w", err)
	}

	vars := terraformVars(name)

	err = terraformCommand(t.WorkingDir, TerraformApply, vars).Run(ctx, t.StdOut, t.StdErr)
	if err != nil {
		return fmt.Errorf("failed to apply terraform: %w", err)
	}

	return nil
}

func (t Terraform) Destroy(ctx context.Context, stateBucket string, stateKey string, name string) error {
	backend := backendConfig(stateBucket, stateKey)

	err := terraformInitCommand(t.WorkingDir, backend).Run(ctx, t.StdOut, t.StdErr)
	if err != nil {
		return fmt.Errorf("failed to initialise terraform: %w", err)
	}

	vars := terraformVars(name)

	err = terraformCommand(t.WorkingDir, TerraformDestroy, vars).Run(ctx, t.StdOut, t.StdErr)
	if err != nil {
		return fmt.Errorf("failed to destroy terraform: %w", err)
	}

	return nil
}

func terraformInitCommand(workingDir string, backendConfig []string) runner {
	args := []string{
		"init",
		"-reconfigure",
	}

	args = append(args, backendConfig...)

	return runner{
		Executable: TerraformExecutable,
		WorkingDir: workingDir,
		Arguments:  args,
	}
}

func terraformCommand(workingDir string, apply TerraformCommandType, vars []string) runner {
	args := []string{
		"apply",
		"-auto-approve",
	}

	args = append(args, vars...)

	if !apply {
		args = append(args, "-destroy")
	}

	return runner{
		Executable: TerraformExecutable,
		WorkingDir: workingDir,
		Arguments:  args,
	}
}

func backendConfig(bucket, key string) []string {
	return []string{
		"-backend-config",
		fmt.Sprintf("bucket=%s", bucket),
		"-backend-config",
		fmt.Sprintf("key=%s", key),
	}
}

func terraformVars(name string) []string {
	return []string{
		"-var",
		fmt.Sprintf("name=%s", name),
	}
}

type TerraformContainer struct {
	Client *docker.Client
	Config *docker.Config
}

func (p TerraformContainer) Create(ctx context.Context, stateBucket string, stateKey string, name string) error {
	config := *p.Config
	config.Cmd = []string{
		"infra",
		"create",
		"--stateBucket",
		stateBucket,
		"--stateKey",
		stateKey,
		"--name",
		name,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return fmt.Errorf("failed to create infra: %w", err)
	}

	return nil
}

func (p TerraformContainer) Destroy(ctx context.Context, stateBucket string, stateKey string, name string) error {
	config := *p.Config
	config.Cmd = []string{
		"infra",
		"destroy",
		"--stateBucket",
		stateBucket,
		"--stateKey",
		stateKey,
		"--name",
		name,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return fmt.Errorf("failed to destroy infra: %w", err)
	}

	return nil
}
