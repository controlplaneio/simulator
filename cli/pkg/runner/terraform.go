package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

// Root return the absolute path of the directory containing the terraform scripts
func Root() (string, error) {
	debug("Finding root")
	absPath, err := filepath.Abs(TfDir())
	if err != nil {
		return "", errors.Wrap(err, "Error resolving root")
	}

	return absPath, nil
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

// PrepareTfEnv returns a slice of strins containing key value pairs of environment variables for
// the child process when exec'ing terraform
func PrepareTfEnv() []string {
	return append(os.Environ(), "TF_IS_IN_AUTOMATION=1")
}

// Terraform wraps running terraform as a child process
func Terraform(cmd string) (*string, error) {
	args := PrepareTfArgs(cmd)

	debug("Preparing to run terraform with args: ", args)
	child := exec.Command("terraform", args...)

	child.Env = PrepareTfEnv()

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir, err := Root()
	if err != nil {
		debug("Error finding root")
		return nil, err
	}

	debug("Setting terraform working directory to ", dir)
	child.Dir = dir

	// Copy child stdout to stdout but also into a buffer to be returned
	var buf bytes.Buffer
	tee := io.TeeReader(childOut, &buf)

	debug("Running terraform")
	err = child.Start()
	if err != nil {
		debug("Error starting terraform child process: ", err)
		return nil, err
	}

	io.Copy(os.Stdout, tee)
	io.Copy(os.Stderr, childErr)

	err = child.Wait()
	if err != nil && err.Error() != "exit status 127" {
		debug("Error waiting for terraform child process", err)
		return nil, err
	}

	out := string(buf.Bytes())
	return &out, nil
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
