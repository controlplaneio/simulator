package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PerturbOptions represents the parameters required by the perturb.sh script
type PerturbOptions struct {
	Master       net.IP
	Slaves       []net.IP
	ScenarioName string
}

// MakePerturbOptions takes a TerraformOutput and a path to a scenario and makes a struct of PerturbOptions
func MakePerturbOptions(tfo TerraformOutput, path string) PerturbOptions {
	po := PerturbOptions{
		Master: net.ParseIP(tfo.MasterNodesPrivateIP.Value[0]),
		Slaves: []net.IP{},
	}

	for _, slave := range tfo.ClusterNodesPrivateIP.Value {
		po.Slaves = append(po.Slaves, net.ParseIP(slave))
	}

	// TODO just use the path and get perturb to do the right thing
	// BUG (rem): pertrb should be able to handle an arbitrary path to a scenario dir
	startOfScenarioName := strings.LastIndex(path, "/") + 1

	po.ScenarioName = path[startOfScenarioName:]

	return po
}

// ToArguments converts a PerturbOptions struct into a slice of strings containing the command line options to pass to
// perturb
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

// PerturbRoot return the absolute path of the directory containing the terraform scripts
func PerturbRoot() (string, error) {
	debug("Finding root")
	absPath, err := filepath.Abs("../simulation-scripts")
	if err != nil {
		return "", errors.Wrap(err, "Error resolving root")
	}

	return absPath, nil
}

// Perturb runs the perturb script with the supplied options
func Perturb(po *PerturbOptions) (*string, error) {
	args := po.ToArguments()

	debug("Preparing to run perturb with args: ", args)
	child := exec.Command("./perturb.sh", args...)

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir, err := PerturbRoot()
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
