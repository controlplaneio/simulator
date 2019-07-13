package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Config returns a pointer to string containing the stanzas to add to an ssh config file so that the kubernetes nodes
// are connectable directly via the bastion or an error if the infrastructure has not been created
func Config(logger *zap.SugaredLogger, tfDir, scenarioPath, bucketName string) (*string, error) {
	tfo, err := Status(logger, tfDir, bucketName)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting infrastructure status")
	}

	if !tfo.IsUsable() {
		return nil, errors.Errorf("No infrastructure, please run simulator infra create")
	}

	return tfo.ToSSHConfig()
}

// Attack establishes an SSH connection to the attack container running on the bastion host ready for the user to
// attempt to complete a scenario
func Attack(logger *zap.SugaredLogger, tfDir, bucketName string) error {
	tfo, err := Status(logger, tfDir, bucketName)
	if err != nil {
		return errors.Wrap(err, "Error getting infrastrucutre status")
	}

	bastion := tfo.BastionPublicIP.Value

	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}

	util.Debug("Running key scan")
	err = ssh.EnsureKnownHosts(bastion)
	if err != nil {
		return errors.Wrap(err, "Error writing known hosts")
	}

	util.Debug("Connecting to", bastion)
	ssh.SSH(bastion)
	return nil
}
