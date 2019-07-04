package simulator

import (
	"github.com/pkg/errors"
)

// Config returns a pointer to string containing the stanzas to add to an ssh config file so that the kubernetes nodes
// are connectable directly via the bastion or an error if the infrastructure has not been created
func Config(tfDir, scenarioPath, bucketName string) (*string, error) {
	tfo, err := Status(tfDir, bucketName)
	if err != nil {
		return nil, err
	}

	if !tfo.IsUsable() {
		return nil, errors.Errorf("No infrastructure, please run simulator infra create")
	}

	return tfo.ToSSHConfig()
}
