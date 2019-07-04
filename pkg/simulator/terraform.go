package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
)

const tfStateDir = "/.terraform"

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
func Terraform(wd, cmd string) (*string, error) {
	args := PrepareTfArgs(cmd)
	env := []string{"TF_IS_IN_AUTOMATION=1"}
	return util.Run(wd, env, "terraform", args...)
}

// InitIfNeeded checks if there is a terraform state folder and calls terraform init if not
func InitIfNeeded(tfDir, bucketName string) error {
	stateDir := tfDir + tfStateDir
	hasStateDir, err := util.FileExists(stateDir)
	if err != nil || hasStateDir {
		return errors.Wrapf(err, "Error checking if terraform state dir exists %s", stateDir)
	}

	ip, err := util.DetectPublicIP()
	if err != nil {
		return err
	}
	accessCIDR := *ip + "/32"
	publicKeyPath, err := util.ExpandTilde("~/.ssh/id_rsa.pub")
	if err != nil {
		return err
	}

	publicKey, err := util.Slurp(*publicKeyPath)
	if err != nil {
		return errors.Wrap(err, "Error reading ~/.ssh/id_rsa.pub")
	}

	err = EnsureTfVarsFile(tfDir, *publicKey, accessCIDR, bucketName)

	_, err = Terraform(tfDir, "init")
	if err != nil {
		return errors.Wrap(err, "Error initialising terraform")
	}

	return nil
}

// -#-

// Create runs terraform init, plan, apply to create the necessary infratsructure to run scenarios
func Create(tfDir, bucketName string) error {
	err := InitIfNeeded(tfDir, bucketName)
	if err != nil {
		return err
	}

	_, err = Terraform(tfDir, "plan")
	if err != nil {
		return err
	}

	_, err = Terraform(tfDir, "apply")
	return err
}

// Status calls terraform output to get the state of the infrastruture and parses the output for
// programmatic use
func Status(tfDir, bucketName string) (*TerraformOutput, error) {
	err := InitIfNeeded(tfDir, bucketName)
	if err != nil {
		return nil, err
	}

	out, err := Terraform(tfDir, "output")
	if err != nil {
		return nil, err
	}

	tfo, err := ParseTerraformOutput(*out)
	if err != nil {
		return nil, err
	}

	return tfo, nil
}

// Destroy call terraform destroy to remove the infrastructure
func Destroy(tfDir, bucketName string) error {
	err := InitIfNeeded(tfDir, bucketName)
	if err != nil {
		return err
	}

	_, err = Terraform(tfDir, "destroy")
	return err
}
