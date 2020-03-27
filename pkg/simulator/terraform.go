package simulator

import (
	"fmt"

	"github.com/kubernetes-simulator/simulator/pkg/childminder"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PrepareTfArgs takes a string with the terraform command desired and returns
// a slice of strings containing the complete list of arguments including the
// command to use when exec'ing terraform
func (s *Simulator) PrepareTfArgs(cmd string) []string {
	arguments := []string{cmd}

	if cmd == "output" {
		arguments = append(arguments, "-json")
	}

	if cmd == "init" || cmd == "plan" || cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-input=false")
		arguments = append(arguments, fmt.Sprintf("--var-file=%s/settings/bastion.tfvars", s.TfVarsDir))

	}

	if cmd == "init" {
		providerBucketArg := fmt.Sprintf("-backend-config=bucket=%s", s.BucketName)
		arguments = append(arguments, providerBucketArg)
	}

	if cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "-auto-approve")
	}

	return arguments
}

// Terraform wraps running terraform as a child process
//func Terraform(wd, cmd string, bucket, tfVarsDir string) (*string, error) {
func (s *Simulator) Terraform(cmd string) (*string, error) {
	args := s.PrepareTfArgs(cmd)
	env := []string{"TF_IS_IN_AUTOMATION=1", "TF_INPUT=0"}
	cm := childminder.NewChildMinder(s.Logger, s.TfDir, env, "terraform", args...)
	if cmd == "output" {
		out, _, err := cm.RunSilently()
		return out, err
	}
	return cm.Run()
}

// InitIfNeeded checks the IP address and SSH key and updates the tfvars if
// needed
func (s *Simulator) InitIfNeeded() error {

	var i interface{} = s.DisableIPDetection
	_, isBool := i.(bool)
	if !isBool {
		return errors.New("disable-ip-detection is not a boolean")
	}

	s.Logger.Debug("Terraform.InitIfNeeded() start")
	s.Logger.Info("Ensuring there is a simulator keypair")
	_, err := s.SSHStateProvider.GetSSHKeyPair()
	if err != nil {
		return errors.Wrap(err, "Error ensuring SSH key")
	}

	var accessCIDR string
	if s.DisableIPDetection {
		accessCIDR = ""
	} else {
		s.Logger.Info("Detecting your public IP address")
		ip, err := util.DetectPublicIP()
		if err != nil {
			return errors.Wrap(err, "Error detecting IP address")
		}
		accessCIDR = *ip + "/32"
	}
	s.Logger.Debug("Reading public key")

	keypair, err := s.SSHStateProvider.GetSSHKeyPair()
	if err != nil {
		return errors.Wrap(err, "Error reading SSH keypair")
	}

	s.Logger.WithFields(logrus.Fields{
		"TfDir":      s.TfDir,
		"TfVarsDir":  s.TfVarsDir,
		"PublicKey":  string(keypair.PublicKey),
		"AccessCIDR": accessCIDR,
		"BucketName": s.BucketName,
		"ExtraCIDRs": s.ExtraCIDRs,
	}).Debug("Writing Terraform tfvars file")
	err = EnsureLatestTfVarsFile(s.TfVarsDir, string(keypair.PublicKey), accessCIDR, s.BucketName, s.AttackTag, s.AttackRepo, s.ExtraCIDRs)
	if err != nil {
		return errors.Wrap(err, "Error writing tfvars")
	}

	s.Logger.Info("Running terraform init")
	_, err = s.Terraform("init")
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
	_, err = s.Terraform("plan")
	if err != nil {
		return err
	}

	s.Logger.Info("Running terraform apply")
	_, err = s.Terraform("apply")
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
	out, err := s.Terraform("output")
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

	s.Logger.Info("Running terraform destroy")
	_, err = s.Terraform("destroy")
	return err
}
