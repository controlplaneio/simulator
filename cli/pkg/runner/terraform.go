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
	tfDir      = "../terraform/deployments/AwsSimulatorStandalone"
	tfStateDir = tfDir + "/.terraform"
)

func Root() (string, error) {
	debug("Finding root")
	absPath, err := filepath.Abs(tfDir)
	if err != nil {
		return "", errors.Wrap(err, "Error resolving root")
	}

	return absPath, nil
}

func PrepareArguments(cmd string) []string {
	arguments := []string{cmd}

	if cmd == "output" {
		arguments = append(arguments, "-json")
	}
	if cmd == "plan" || cmd == "apply" || cmd == "destroy" {
		arguments = append(arguments, "--var-file=settings/bastion.tfvars")
		arguments = append(arguments, "-auto-approve")
	}
	if cmd == "init" {
		arguments = append(arguments, "--var-file=settings/bastion.tfvars")
	}

	return arguments

}

func Terraform(cmd string) (*string, error) {
	arguments := PrepareArguments(cmd)

	debug("Preparing to run terraform with args: ", arguments)
	child := exec.Command("terraform", arguments...)

	// Tell terraform it is being automated
	child.Env = append(os.Environ(), "TF_IS_IN_AUTOMATION=1")

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
		debug("Error starting terraform child process")
		return nil, err
	}

	io.Copy(os.Stdout, tee)
	io.Copy(os.Stderr, childErr)

	err = child.Wait()
	if err != nil {
		debug("Error waiting for terraform child process")
		return nil, err
	}

	out := string(buf.Bytes())
	return &out, nil
}

func InitIfNeeded() error {
	hasStateDir, err := exists(tfStateDir)
	if err != nil {
		return errors.Wrapf(err, "Error checking if terraform state dir exists %s", tfStateDir)
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
