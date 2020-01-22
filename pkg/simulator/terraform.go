package simulator

import (
	"fmt"

	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
)

// PrepareTfArgs takes a string with the terraform command desired and returns
// a slice of strings containing the complete list of arguments including the
// command to use when exec'ing terraform
func PrepareTfArgs(cmd string, bucket, tfVarsDir string) []string {
	arguments := []string{cmd}

	if cmd == "output" {
		arguments = append(arguments, "-json")
	}

	if cmd == "init" || cmd == "plan" || cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-input=false")
		arguments = append(arguments, "--var-file=" + tfVarsDir + "/settings/bastion.tfvars")

	}

	if cmd == "init" {
		providerBucketArg := fmt.Sprintf("-backend-config=bucket=%s", bucket)
		arguments = append(arguments, providerBucketArg)
	}

	if cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-auto-approve")
	}

	return arguments
}

// Terraform wraps running terraform as a child process
func Terraform(wd, cmd string, bucket, tfVarsDir string) (*string, error) {
	args := PrepareTfArgs(cmd, bucket, tfVarsDir)
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
func (s *Simulator) InitIfNeeded() error {
	s.Logger.Debug("Terraform.InitIfNeeded() start")
	s.Logger.Info("Ensuring there is a simulator keypair")
	_, err := ssh.EnsureKey()
	if err != nil {
		return errors.Wrap(err, "Error ensuring SSH key")
	}

	s.Logger.Info("Detecting your public IP address")
	ip, err := util.DetectPublicIP()
	if err != nil {
		return errors.Wrap(err, "Error detecting IP address")
	}
	accessCIDR := *ip + "/32"

	s.Logger.Debug("Reading public key")
	publickey, err := ssh.PublicKey()
	if err != nil {
		return errors.Wrap(err, "Error reading public key")
	}

	s.Logger.Debugf("terraform Directory: %s", s.TfDir)
	s.Logger.Debugf("terraform vars Directory: %s", s.TfVarsDir)
	s.Logger.Debugf("Public Key:\n%s", publickey)
	s.Logger.Debugf("Access CIDR: %s", accessCIDR)
	s.Logger.Debugf("Remote State Bucket Name: %s", s.BucketName)
	s.Logger.Debug("Writing terraform tfvars")
	err = EnsureLatestTfVarsFile(s.TfVarsDir, *publickey, accessCIDR, s.BucketName, s.AttackTag)
	if err != nil {
		return errors.Wrap(err, "Error writing tfvars")
	}

	s.Logger.Info("Running terraform init")
	_, err = Terraform(s.TfDir, "init", s.BucketName, s.TfVarsDir)
	if err != nil {
		return errors.Wrap(err, "Error initialising terraform")
	}

	return nil
}

// -#-

// Create runs terraform init, plan, apply to create the necessary
// infrastructure to run scenarios
func (s *Simulator) Create() error {

	err := s.InitIfNeeded()

	if err != nil {
		return err
	}
	
	s.Logger.Info("Running terraform plan")
	_, err = Terraform(s.TfDir, "plan", s.BucketName, s.TfVarsDir)
	if err != nil {
		return err
	}

	s.Logger.Info("Running terraform apply")
	_, err = Terraform(s.TfDir, "apply", s.BucketName, s.TfVarsDir)
	return err
}

// Status calls terraform output to get the state of the infrastruture and
// parses the output for programmatic use
func (s *Simulator) Status() (*TerraformOutput, error) {
	//err := s.InitIfNeeded()
	err := s.InitIfNeeded()
	if err != nil {
		return nil, errors.Wrap(err, "Error initialising")
	}

	s.Logger.Info("Running terraform output")
	out, err := Terraform(s.TfDir, "output", s.BucketName, s.TfVarsDir)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting terraform outputs")
	}

	s.Logger.Debug(out)

	s.Logger.Debug("Parsing terraform output")
	tfo, err := ParseTerraformOutput(*out)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing terraform outputs")
	}

	return tfo, nil
}

// Destroy call terraform destroy to remove the infrastructure
func (s *Simulator) Destroy() error {
	err := s.InitIfNeeded()
	if err != nil {
		return errors.Wrap(err, "Error initialising")
	}

	s.Logger.Info("Running terrraform destroy")
	_, err = Terraform(s.TfDir, "destroy", s.BucketName, s.TfVarsDir)
	return err
}
