package ssh

import (
	"github.com/kubernetes-simulator/simulator/pkg/childminder"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"os"
)

// EnsureKnownHosts scans the bastion host for its SSH keys and writes them to
// a custom known hosts location for promptless interaction with the
// infrastructure.  Return an error if any occurred
func EnsureKnownHosts(bastion string) error {
	hostkeys, err := KeyScan(bastion)
	if err != nil {
		return errors.Wrap(err, "Error running ssh-keyscan")
	}

	abspath, err := util.ExpandTilde(KnownHostsPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving SSH known hosts path")
	}

	err = util.OverwriteFile(*abspath, "# "+bastion+"\n"+*hostkeys)
	if err != nil {
		return errors.Wrap(err, "Error writing SSH known hosts file")
	}

	return nil
}

// KeyScan runs ssh-keyscan silently against the provided bastion address. It
// returns a pointer to a string containing its buffered stdout or an error if
// any occurred
func KeyScan(bastion string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting process working directory")
	}
	cm := childminder.NewChildMinder(nil, wd, os.Environ(), "ssh-keyscan", "-H", bastion)
	out, _, err := cm.RunSilently()
	return out, err
}
