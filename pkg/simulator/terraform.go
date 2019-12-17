package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"fmt"
)

// PrepareTfArgs takes a string with the terraform command desired and returns
// a slice of strings containing the complete list of arguments including the
// command to use when exec'ing terraform
func PrepareTfArgs(cmd string, bucket string) []string {
	arguments := []string{cmd}

	if cmd == "output" {
		arguments = append(arguments, "-json")
	}

	if cmd == "init" || cmd == "plan" || cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-input=false")
		arguments = append(arguments, "--var-file=settings/bastion.tfvars")

	}

	if cmd == "init" {
		providerBucketFlag := fmt.Sprintf("-backend-config=bucket=%s", bucket)
		arguments = append(arguments, providerBucketFlag)
	}

	if cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-auto-approve")
	}

	return arguments
}

// Terraform wraps running terraform as a child process
func Terraform(wd, cmd string, bucket string) (*string, error) {
	args := PrepareTfArgs(cmd, bucket)
	env := []string{"TF_IS_IN_AUTOMATION=1", "TF_INPUT=0"}
	if cmd == "output" {
		// TODO: (rem) deal with non-empty stderr?
		out, _, err := util.RunSilently(wd, env, "terraform", args...)
		return out, err
	}

	return util.Run(wd, env, "terraform", args...)
}

// InitIfNeeded checks the IP address and SSH key and updates the tfvars if
// needed
func InitIfNeeded(logger *zap.SugaredLogger, tfDir, bucket, attackTag string) error {
	logger.Debug("Terraform.InitIfNeeded() start")

	logger.Info("Ensuring there is a simulator keypair")
	_, err := ssh.EnsureKey()
	if err != nil {
		return errors.Wrap(err, "Error ensuring SSH key")
	}

	logger.Info("Detecting your public IP address")
	ip, err := util.DetectPublicIP()
	if err != nil {
		return errors.Wrap(err, "Error detecting IP address")
	}
	accessCIDR := *ip + "/32"

	logger.Debug("Reading public key")
	publickey, err := ssh.PublicKey()
	if err != nil {
		return errors.Wrap(err, "Error reading public key")
	}

	logger.Debugf("terraform Directory: %s", tfDir)
	logger.Debugf("Public Key:\n%s", publickey)
	logger.Debugf("Access CIDR: %s", accessCIDR)
	logger.Debugf("Remote State Bucket Name: %s", bucket)
	logger.Debug("Writing terraform tfvars")
	err = EnsureLatestTfVarsFile(tfDir, *publickey, accessCIDR, bucket, attackTag)
	if err != nil {
		return errors.Wrap(err, "Error writing tfvars")
	}

	logger.Info("Running terraform init")
	_, err = Terraform(tfDir, "init", bucket)
	if err != nil {
		return errors.Wrap(err, "Error initialising terraform")
	}

	return nil
}

// -#-

// Create runs terraform init, plan, apply to create the necessary
// infrastructure to run scenarios
func Create(logger *zap.SugaredLogger, tfDir, bucket, attackTag string) error {
	err := InitIfNeeded(logger, tfDir, bucket, attackTag)

	if err != nil {
		return err
	}

	logger.Info("Running terraform plan")
	_, err = Terraform(tfDir, "plan", bucket)
	if err != nil {
		return err
	}

	logger.Info("Running terraform apply")
	_, err = Terraform(tfDir, "apply", bucket)
	return err
}

// Status calls terraform output to get the state of the infrastruture and
// parses the output for programmatic use
func Status(logger *zap.SugaredLogger, tfDir, bucket, attackTag string) (*TerraformOutput, error) {
	err := InitIfNeeded(logger, tfDir, bucket, attackTag)
	if err != nil {
		return nil, errors.Wrap(err, "Error initialising")
	}

	logger.Info("Running terraform output")
	out, err := Terraform(tfDir, "output", bucket)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting terraform outputs")
	}

	logger.Debug(out)

	logger.Debug("Parsing terraform output")
	tfo, err := ParseTerraformOutput(*out)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing terraform outputs")
	}

	return tfo, nil
}

// Destroy call terraform destroy to remove the infrastructure
func Destroy(logger *zap.SugaredLogger, tfDir, bucket, attackTag string) error {
	err := InitIfNeeded(logger, tfDir, bucket, attackTag)
	if err != nil {
		return errors.Wrap(err, "Error initialising")
	}

	logger.Info("Running terrraform destroy")
	_, err = Terraform(tfDir, "destroy", bucket)
	return err
}
