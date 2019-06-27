package runner

import (
	"github.com/pkg/errors"
	"os"
)

const (
	tfDirEnvVar  = "SIMULATOR_TF_DIR"
	defaultTfDir = "../terraform/deployments/AwsSimulatorStandalone"
	tfStateDir   = "/.terraform"
)

// TfDir reads the Terraform directory from the environment variable `SIMULATOR_TF_DIR`
// or uses a default value of `../terraform/deployments/AwsSimulatorStandalone`
func TfDir() string {
	var d = os.Getenv(tfDirEnvVar)
	if d == "" {
		d = defaultTfDir
	}

	return d
}

// PrepareTfArgs takes a string with the terraform command desired and returns a slice of strings
// containing the complete list of arguments including the command to use when exec'ing terraform
func PrepareTfArgs(cmd string) []string {
	arguments := []string{cmd}

	if cmd == "output" {
		arguments = append(arguments, "-json")
	}
	if cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "--var-file=settings/bastion.tfvars")
		arguments = append(arguments, "-auto-approve")
	}

	if cmd == "init" || cmd == "plan" {
		arguments = append(arguments, "--var-file=settings/bastion.tfvars")
	}

	return arguments
}

// Terraform wraps running terraform as a child process
func Terraform(cmd string) (*string, error) {
	args := PrepareTfArgs(cmd)
	env := []string{"TF_IS_IN_AUTOMATION=1"}
	wd := TfDir()
	return Run(wd, env, "terraform", args...)
}

// InitIfNeeded checks if there is a terraform state folder and calls terraform init if not
func InitIfNeeded() error {
	stateDir := TfDir() + tfStateDir
	hasStateDir, err := exists(stateDir)
	if err != nil {
		return errors.Wrapf(err, "Error checking if terraform state dir exists %s", stateDir)
	}

	if hasStateDir {
		return nil
	}

	_, err = Terraform("init")
	if err != nil {
		return errors.Wrap(err, "Error initialising terraform")
	}

	return nil
}

// -#-

// Create runs terraform init, plan, apply to create the necessary infratsructure to run scenarios
func Create() error {
	err := InitIfNeeded()
	if err != nil {
		return err
	}

	_, err = Terraform("plan")
	if err != nil {
		return err
	}

	_, err = Terraform("apply")
	return err
}

// Status calls terraform output to get the state of the infrastruture and parses the output for
// programmatic use
func Status() (*TerraformOutput, error) {
	err := InitIfNeeded()
	if err != nil {
		return nil, err
	}

	out, err := Terraform("output")
	tfo, err := ParseTerraformOutput(*out)
	if err != nil {
		return nil, err
	}

	return tfo, nil
}

// Destroy call terraform destroy to remove the infrastructure
func Destroy() error {
	err := InitIfNeeded()
	if err != nil {
		return err
	}

	_, err = Terraform("destroy")
	return err
}
