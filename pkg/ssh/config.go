package ssh

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
)

// EnsureSSHConfig writes the SSH config file for simulator that is needed for the `perturb.sh` scripts to
// run succesfully via the bastion
func EnsureSSHConfig(cfg string) error {
	abspath, err := util.ExpandTilde(SSHConfigPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving SSH config path")
	}

	err = util.OverwriteFile(*abspath, cfg)
	if err != nil {
		return errors.Wrap(err, "Error overwriting SSH config")
	}

	return nil
}
