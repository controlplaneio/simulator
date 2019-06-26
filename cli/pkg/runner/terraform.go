package runner

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Root() (string, error) {
	absPath, err := filepath.Abs("../terraform/deployments/AWSSimulatorStandalone")
	if err != nil {
		return "", errors.Wrap(err, "Error resolving root")
	}

	return absPath, nil
}

func Terraform(cmd string) error {
	child := exec.Command("terraform", cmd, "--var-file=settings/bastion.tfvars")

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()

	tfDir, err := Root()
	if err != nil {
		return err
	}

	child.Dir = tfDir

	child.Start()

	io.Copy(os.Stdout, childOut)
	io.Copy(os.Stderr, childErr)
	childIn.Close()

	return child.Wait()
}

func Create() error {
	err := Terraform("init")
	if err != nil {
		return err
	}

	err = Terraform("plan")
	if err != nil {
		return err
	}

	return Terraform("apply")
}

func Destroy() error {
	err := Terraform("init")
	if err != nil {
		return err
	}

	return Terraform("destroy")
}
