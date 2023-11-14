package controlplane

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/controlplaneio/simulator/v2/controlplane/aws"
	"github.com/controlplaneio/simulator/v2/controlplane/commands"
)

const (
	SimulatorDir = "/simulator"
	Ansible      = "ansible"
	Packer       = "packer"
	Terraform    = "terraform"
)

var (
	AWSDir             = "/home/ubuntu/.aws"
	AdminSSHBundleDir  = filepath.Join(SimulatorDir, "admin")
	PlayerSSHBundleDir = filepath.Join(SimulatorDir, "player")

	AnsibleDir         = filepath.Join(SimulatorDir, Ansible)
	AnsiblePlaybookDir = filepath.Join(AnsibleDir, "playbooks")

	PackerTemplateDir = filepath.Join(SimulatorDir, Packer)

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

type simulator struct{}

func (s simulator) CreateBucket(ctx context.Context, name string) error {
	//nolint:wrapcheck
	return aws.CreateBucket(ctx, name)
}

func (s simulator) BuildImage(ctx context.Context, name string) error {
	err := commands.PackerInitCommand(PackerTemplateDir, name).Run(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	//nolint:wrapcheck
	return commands.PackerBuildCommand(PackerTemplateDir, name).Run(ctx)
}

func (s simulator) CreateInfrastructure(ctx context.Context, bucket, key, name string) error {
	err := commands.TerraformInitCommand(TerraformWorkspaceDir, backendConfig(bucket, key)).Run(ctx)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	//nolint:wrapcheck
	return commands.TerraformCommand(TerraformWorkspaceDir, commands.TerraformApply, terraformVars(name, bucket)).
		Run(ctx)
}

func (s simulator) DestroyInfrastructure(ctx context.Context, bucket, key, name string) error {
	err := commands.TerraformInitCommand(TerraformWorkspaceDir, backendConfig(bucket, key)).Run(ctx)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	//nolint:wrapcheck
	return commands.TerraformCommand(TerraformWorkspaceDir, commands.TerraformDestroy, terraformVars(name, bucket)).
		Run(ctx)
}

func (s simulator) InstallScenario(ctx context.Context, name string) error {
	//nolint:wrapcheck
	return commands.AnsiblePlaybookCommand(AdminSSHBundleDir, AnsiblePlaybookDir, name).Run(ctx)
}

func (s simulator) UninstallScenario(ctx context.Context, name string) error {
	//nolint:wrapcheck
	return commands.AnsiblePlaybookCommand(AdminSSHBundleDir, AnsiblePlaybookDir, name, "state=absent").
		Run(ctx)
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
		fmt.Sprintf("admin_ssh_bundle_dir=%s", AdminSSHBundleDir),
		"-var",
		fmt.Sprintf("player_ssh_bundle_dir=%s", PlayerSSHBundleDir),
	}
}
