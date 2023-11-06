package controlplane

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/controlplaneio/simulator/controlplane/aws"
	"github.com/controlplaneio/simulator/controlplane/commands"
)

const (
	SimulatorDir = "/simulator"
	Home         = "config"
	Scenarios    = "scenarios"
	Packer       = "packer"
	Terraform    = "terraform"
)

var (
	HomeDir         = filepath.Join(SimulatorDir, Home)
	AdminConfigDir  = filepath.Join(HomeDir, "admin")
	PlayerConfigDir = filepath.Join(HomeDir, "player")

	AnsibleDir                = filepath.Join(SimulatorDir, Scenarios)
	AnsiblePlaybookDir string = filepath.Join(AnsibleDir, "playbooks")

	PackerTemplateDir string = filepath.Join(SimulatorDir, Packer)

	TerraformDir          = filepath.Join(SimulatorDir, Terraform)
	TerraformWorkspaceDir = filepath.Join(TerraformDir, "workspaces/simulator")
)

func New() Simulator {
	return simulator{}
}

type Simulator interface {
	// CreateBucket creates the S3 bucket used to store Simulator state.
	CreateBucket(ctx context.Context, name string) error
	// BuildImage runs Packer to create the specified AMI.
	BuildImage(ctx context.Context, name string) error
	// CreateInfrastructure runs Terraform to create the Simulator infrastructure.
	CreateInfrastructure(ctx context.Context, bucket string, key string, name string) error
	// DestroyInfrastructure runs Terraform to destroy the Simulator infrastructure.
	DestroyInfrastructure(ctx context.Context, bucket string, key string, name string) error
	// InstallScenario installs a scenario on the Simulator infrastructure.
	InstallScenario(ctx context.Context, name string) error
	// UninstallScenario uninstalls a scenario from the Simulator infrastructure.
	UninstallScenario(ctx context.Context, name string) error
}

type simulator struct {
}

func (s simulator) CreateBucket(ctx context.Context, name string) error {
	slog.Debug("simulator init", "name", name)

	return aws.CreateBucket(ctx, name)
}

func (s simulator) BuildImage(ctx context.Context, name string) error {
	slog.Debug("simulator build", "image", name)

	err := commands.PackerInitCommand(PackerTemplateDir, string(name)).Run(ctx)
	if err != nil {
		return err
	}

	return commands.PackerBuildCommand(PackerTemplateDir, string(name)).Run(ctx)
}

// TODO: add state path and config bucket and path to support Kubesim
func (s simulator) CreateInfrastructure(ctx context.Context, bucket, key, name string) error {
	slog.Debug("simulator create infrastructure", "bucket", bucket, "key", key, "name", name)

	err := commands.TerraformInitCommand(TerraformWorkspaceDir, backendConfig(bucket, key)).Run(ctx)
	if err != nil {
		return err
	}

	return commands.TerraformCommand(TerraformWorkspaceDir, commands.TerraformApply, terraformVars(name, bucket)).Run(ctx)
}

func (s simulator) DestroyInfrastructure(ctx context.Context, bucket, key, name string) error {
	slog.Debug("simulator destroy infrastructure", "bucket", bucket, "key", key, "name", name)

	err := commands.TerraformInitCommand(TerraformWorkspaceDir, backendConfig(bucket, key)).Run(ctx)
	if err != nil {
		return err
	}

	return commands.TerraformCommand(TerraformWorkspaceDir, commands.TerraformDestroy, terraformVars(name, bucket)).Run(ctx)
}

func (s simulator) InstallScenario(ctx context.Context, name string) error {
	slog.Debug("simulator install", "scenario", name)

	return commands.AnsiblePlaybookCommand(AdminConfigDir, AnsiblePlaybookDir, name).Run(ctx)
}

func (s simulator) UninstallScenario(ctx context.Context, name string) error {
	slog.Debug("simulator uninstall", "scenario", name)

	return commands.AnsiblePlaybookCommand(AdminConfigDir, AnsiblePlaybookDir, name, "state=absent").Run(ctx)
}

func backendConfig(bucket, key string) []string {
	return []string{
		"-backend-config",
		fmt.Sprintf("bucket=%s", bucket),
		"-backend-config",
		fmt.Sprintf("key=%s", key),
	}
}

func terraformVars(name, bucket string) []string {
	return []string{
		"-var",
		fmt.Sprintf("name=%s", name),
		"-var",
		fmt.Sprintf("bucket=%s", bucket),
		"-var",
		fmt.Sprintf("admin_config_dir=%s", AdminConfigDir),
		"-var",
		fmt.Sprintf("player_config_dir=%s", PlayerConfigDir),
	}
}
