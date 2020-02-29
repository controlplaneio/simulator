package simulator

import (
	"net"
	"strings"

	"github.com/controlplaneio/simulator-standalone/pkg/childminder"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/sirupsen/logrus"
)

// PerturbOptions represents the parameters required by the perturb.sh script
type PerturbOptions struct {
	bastion      net.IP
	master       net.IP
	slaves       []net.IP
	scenarioName string
	Force        bool
}

// MakePerturbOptions takes a TerraformOutput and a path to a scenario and
// makes a struct of PerturbOptions
func (po *PerturbOptions) MakePerturbOptions(tfo TerraformOutput, path string) {
	po.master = net.ParseIP(tfo.MasterNodesPrivateIP.Value[0])
	po.slaves = []net.IP{}

	for _, slave := range tfo.ClusterNodesPrivateIP.Value {
		po.slaves = append(po.slaves, net.ParseIP(slave))
	}

	// TODO: (rem) just use the path and get perturb to do the right thing
	// BUG: (rem) perturb should be able to handle an arbitrary path to a scenario dir
	startOfScenarioName := strings.LastIndex(path, "/") + 1

	po.bastion = net.ParseIP(tfo.BastionPublicIP.Value)
	po.scenarioName = path[startOfScenarioName:]
}

// ToArguments converts a PerturbOptions struct into a slice of strings
// containing the command line options to pass to perturb
func (po *PerturbOptions) ToArguments() []string {
	arguments := []string{"--master", po.master.String()}
	arguments = append(arguments, "--bastion")
	arguments = append(arguments, po.bastion.String())
	arguments = append(arguments, "--nodes")
	slaves := ""
	for index, slave := range po.slaves {
		slaves += slave.String()
		if index < len(po.slaves)-1 {
			slaves += ","
		}
	}
	arguments = append(arguments, slaves)
	if po.Force {
		arguments = append(arguments, "--force")
	}
	arguments = append(arguments, po.scenarioName)

	return arguments
}

func (po *PerturbOptions) String() string {
	return strings.TrimSpace(strings.Join(po.ToArguments(), " "))
}

const (
	perturbPathEnvVar  = "SIMULATOR_SCENARIOS_DIR"
	defaultPerturbPath = "./simulation-scripts/"
)

// Perturb runs the perturb script with the supplied options
func Perturb(po *PerturbOptions, logger *logrus.Logger) (*string, error) {
	args := po.ToArguments()
	env := []string{}
	wd := util.EnvOrDefault(perturbPathEnvVar, defaultPerturbPath)
	cm := childminder.NewChildMinder(logger, wd, env, "./perturb.sh", args...)
	return cm.Run()
}
