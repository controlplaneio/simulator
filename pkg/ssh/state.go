package ssh

import (
	"fmt"
	"github.com/kubernetes-simulator/simulator/pkg/childminder"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// StateProvider provides methods for storing or retrieving state about a user and
// their cluster
type StateProvider interface {
	GetSSHKeyPair() (*KeyPair, error)
	SaveSSHConfig(config string) error
	GetSSHConfig() (*string, error)
}

// LocalStateProvider is the default State provider and persists all state into the
// local ~/.kubesim directory
type LocalStateProvider struct{}

// GetSSHKeyPair returns an existing SSH keypair or creates one locally
func (ls LocalStateProvider) GetSSHKeyPair() (*KeyPair, error) {
	abspath, err := util.ExpandTilde(PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error resolving key path")
	}

	exists, err := util.FileExists(*abspath)
	if err != nil {
		return nil, errors.Wrap(err, "Error checking if key already exists")
	}

	// key already exists return it
	if exists {
		return ls.getSSHKeyPair()
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting process working directory")
	}

	cm := childminder.NewChildMinder(nil, wd, os.Environ(), "ssh-keygen",
		"-f", *abspath, "-t", "rsa", "-C",
		"simulator-key", "-N", "")
	stdout, stderr, err := cm.RunSilently()
	if err != nil {
		fmt.Println(*stdout)
		fmt.Println(*stderr)
		return nil, errors.Wrap(err, "Error generating keypair")
	}

	return ls.getSSHKeyPair()
}

// GetSSHKeyPair retieves
func (ls LocalStateProvider) getSSHKeyPair() (*KeyPair, error) {
	publicKeyPath, err := util.ExpandTilde(PublicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", PublicKeyPath)
	}

	publickey, err := util.Slurp(*publicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", PublicKeyPath)
	}

	privateKeyPath, err := util.ExpandTilde(PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", PrivateKeyPath)
	}

	privatekey, err := util.Slurp(*privateKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", PublicKeyPath)
	}
	ret := KeyPair{
		PublicKey:  PublicKey(strings.Trim(*publickey, "\n")),
		PrivateKey: PrivateKey(*privatekey),
	}
	return &ret, nil
}

// SaveSSHConfig saves the config supplied to the local ~/.ssh directory
func (ls LocalStateProvider) SaveSSHConfig(config string) error {
	abspath, err := util.ExpandTilde(ConfigPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving SSH config path")
	}

	err = util.OverwriteFile(*abspath, config)
	if err != nil {
		return errors.Wrap(err, "Error overwriting SSH config")
	}

	return nil
}

// GetSSHConfig reads the config from the local ~/.ssh directory
func (ls LocalStateProvider) GetSSHConfig() (*string, error) {
	abspath, err := util.ExpandTilde(ConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", ConfigPath)
	}

	config, err := util.Slurp(*abspath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", ConfigPath)
	}
	return config, nil
}
