package simulator

import (
	"fmt"
	"github.com/pkg/errors"
)

// Config returns a pointer to string containing the stanzas to add to an ssh config file so that the kubernetes nodes
// are connectable directly via the bastion or an error if the infrastructure has not been created
func Config(tfDir, scenarioPath string) (*string, error) {
	tfo, err := Status(tfDir)
	if err != nil {
		return nil, err
	}

	if !tfo.IsUsable() {
		return nil, errors.Errorf("No infrastructure, please run simulator infra create:\n %#v", tfo)
	}

	po := MakePerturbOptions(*tfo, scenarioPath)
	fmt.Println("Converted usable terraform output into perturb options")
	fmt.Printf("%#v", po)
	return tfo.ToSSHConfig()
}
