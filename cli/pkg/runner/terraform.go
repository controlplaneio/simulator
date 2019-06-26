package runner

import (
	"bytes"
	"fmt"
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

// Reads the Terraform dir from the environment variable `SIMULATOR_TF_DIR`
// or uses a default value of `../terraform/deployments/AwsSimulatorStandalone`
func TfDir() string {
	var d = os.Getenv(tfDirEnvVar)
	if d == "" {
		d = defaultTfDir
	}

	return d
}
func Root() (string, error) {
	debug("Finding root")
	absPath, err := filepath.Abs(TfDir())
	if err != nil {
		return "", errors.Wrap(err, "Error resolving root")
	}

	return absPath, nil
}

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

func PrepareTfEnv() []string {
	return append(os.Environ(), "TF_IS_IN_AUTOMATION=1")
}

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

	tfDir, err := Root()
	if err != nil {
		debug("Error finding root")
		return nil, err
	}

	debug("Setting terraform working directory to ", tfDir)
	child.Dir = tfDir

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

func Status() error {
	err := InitIfNeeded()
	if err != nil {
		return err
	}

	out, err := Terraform("output")
	tfOutput, err := ParseTerraformOutput(*out)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", tfOutput)

	return err
}

func Destroy() error {
	err := InitIfNeeded()
	if err != nil {
		return err
	}

	_, err = Terraform("destroy")
	return err
}
