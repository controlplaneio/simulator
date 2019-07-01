package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"net"
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

	// TODO: (rem) just use the path and get perturb to do the right thing
	// BUG: (rem) pertrb should be able to handle an arbitrary path to a scenario dir
	startOfScenarioName := strings.LastIndex(path, "/") + 1

	po.ScenarioName = path[startOfScenarioName:]

	return po
}

// ToArguments converts a PerturbOptions struct into a slice of strings containing the command line options to pass to
// perturb
func (po *PerturbOptions) ToArguments() []string {
	arguments := []string{"--master", po.Master.String()}
	arguments = append(arguments, "--slaves")
	slaves := ""
	for index, slave := range po.Slaves {
		slaves += slave.String()
		if index < len(po.Slaves)-1 {
			slaves += ","
		}
	}
	arguments = append(arguments, slaves)

	return arguments
}

func (po *PerturbOptions) String() string {
	return strings.Join(po.ToArguments(), " ")
}

const (
	perturbPathEnvVar  = "SIMULATOR_MANIFEST_PATH"
	defaultPerturbPath = "./simulation-scripts/"
)

// Perturb runs the perturb script with the supplied options
func Perturb(po *PerturbOptions) (*string, error) {
	args := po.ToArguments()
	env := []string{}
	wd := util.EnvOrDefault(perturbPathEnvVar, defaultPerturbPath)
	return util.Run(wd, env, "./perturb.sh", args...)
}
