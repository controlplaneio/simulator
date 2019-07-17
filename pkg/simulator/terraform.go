package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

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
	if cmd == "output" {
		// TODO: (rem) deal with non-empty stderr?
		out, _, err := util.RunSilently(wd, env, "terraform", args...)
		return out, err
	}

	return util.Run(wd, env, "terraform", args...)
}

// InitIfNeeded checks the IP address and SSH key and updates the tfvars if needed
func InitIfNeeded(logger *zap.SugaredLogger, tfDir, bucketName string) error {
	logger.Debug("Terraform.InitIfNeeded() start")

	_, err := ssh.EnsureKey()
	if err != nil {
		return errors.Wrap(err, "Error ensuring SSH key")
	}

	ip, err := util.DetectPublicIP()
	if err != nil {
		return errors.Wrap(err, "Error detecting IP address")
	}
	accessCIDR := *ip + "/32"

	publickey, err := ssh.PublicKey()
	if err != nil {
		return errors.Wrap(err, "Error reading public key")
	}

	err = EnsureLatestTfVarsFile(tfDir, *publickey, accessCIDR, bucketName)
	if err != nil {
		return errors.Wrap(err, "Error writing tfvars")
	}

	_, err = Terraform(tfDir, "init")
	if err != nil {
		return errors.Wrap(err, "Error initialising terraform")
	}

	return nil
}

// -#-

// Create runs terraform init, plan, apply to create the necessary infratsructure to run scenarios
func Create(logger *zap.SugaredLogger, tfDir, bucketName string) error {
	err := InitIfNeeded(logger, tfDir, bucketName)

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

// Status calls terraform output to get the state of the infrastruture and parses the output for programmatic use
func Status(logger *zap.SugaredLogger, tfDir, bucketName string) (*TerraformOutput, error) {
	err := InitIfNeeded(logger, tfDir, bucketName)
	if err != nil {
		return nil, errors.Wrap(err, "Error initialising")
	}

	out, err := Terraform(tfDir, "output")
	if err != nil {
		return nil, errors.Wrap(err, "Error getting terraform outputs")
	}

	tfo, err := ParseTerraformOutput(*out)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing terraform outputs")
	}

	return tfo, nil
}

// Destroy call terraform destroy to remove the infrastructure
func Destroy(logger *zap.SugaredLogger, tfDir, bucketName string) error {
	err := InitIfNeeded(logger, tfDir, bucketName)
	if err != nil {
		return errors.Wrap(err, "Error initialising")
	}

	_, err = Terraform(tfDir, "destroy")
	return err
}
