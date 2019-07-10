package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"os"
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

func Attack(tfDir, bucketName string) error {
	tfo, err := Status(tfDir, bucketName)
	if err != nil {
		return err
	}

	bastion := tfo.BastionPublicIP.Value

	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}

	util.Debug("Running key scan")
	hostkeys, err := keyScan(bastion)
	if err != nil {
		return err
	}

	knownhosts, err := util.ExpandTilde("~/.ssh/known_hosts")
	if err != nil {
		return err
	}
	written, err := util.EnsureFile(*knownhosts, "# "+bastion+"\n"+*hostkeys)
	if err != nil {
		return err
	}

	if !written {
		util.Debug("Did not write ", *hostkeys)
	}

	util.Debug("Connecting to", bastion)
	util.SSH(bastion)
	return nil
}

func keyScan(bastion string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := util.RunSilently(wd, os.Environ(), "ssh-keyscan", "-H", bastion)
	return out, err
}
