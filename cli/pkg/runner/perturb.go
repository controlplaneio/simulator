package runner

import (
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type PerturbOptions struct {
	Master net.IP
	Slaves []net.IP
}

func (po *PerturbOptions) ToArguments() []string {
	arguments := []string{"--master", po.Master.String()}
	arguments = append(arguments, "--slaves")
	for index, slave := range po.Slaves {
		s := slave.String()
		if index < len(po.Slaves)-1 {
			s += ","
		}

		arguments = append(arguments, s)
	}
	return arguments
}

func (po *PerturbOptions) String() string {
	return strings.Join(po.ToArguments(), " ")
}

func Perturb(po *PerturbOptions) error {
	child := exec.Command("../simulation-scripts/perturb.sh", po.ToArguments()...)

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
