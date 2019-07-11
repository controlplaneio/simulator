package ssh

import (
	"encoding/base64"
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

// Base64PrivateKey returns a pointer to a string containing the base64 encoded private key or an error
func Base64PrivateKey(name string) (*string, error) {
	keypath, err := util.ExpandTilde(name)
	if err != nil {
		return nil, err
	}

	keymaterial, err := util.Slurp(*keypath)
	if err != nil {
		return nil, err
	}

	encodedkey := base64.StdEncoding.EncodeToString([]byte(*keymaterial))

	return &encodedkey, nil
}

// PublicKey reads the public key and return a pointer to a string with its contents or any error
func PublicKey() (*string, error) {
	publicKeyPath, err := util.ExpandTilde(PublicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", PublicKeyPath)
	}

	publickey, err := util.Slurp(*publicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", PublicKeyPath)
	}

	ret := strings.Trim(*publickey, "\n")
	return &ret, nil
}

// UpdateKnownHosts scans the bastion host for its SSH keys and writes them to a custom known hosts location for
// promptless interaction with the infrastructure.  Return an error if any occurred
func UpdateKnownHosts(bastion string) error {
	hostkeys, err := KeyScan(bastion)
	if err != nil {
		return errors.Wrap(err, "Error running ssh-keyscan")
	}

	abspath, err := util.ExpandTilde(SSHKnownHostsPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving SSH known hosts path")
	}

	err = util.OverwriteFile(*abspath, "# "+bastion+"\n"+*hostkeys)
	if err != nil {
		return errors.Wrap(err, "Error writing SSH known hosts file")
	}

	return nil
}

// KeyScan runs ssh-keyscan silently against the provided bastion address. It returns a pointer to a string containing
// its buffered stdout or an error if any occurred
func KeyScan(bastion string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := util.RunSilently(wd, os.Environ(), "ssh-keyscan", "-H", bastion)
	return out, err
}

// GenerateKey runs ssh-keygen silently to create an SSH key with the same provided using preconfigured settings
// It returns a pointer to a string containing the buffered stdout or an error if any occurred
func GenerateKey(privatekeypath string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, stderr, err := util.RunSilently(wd, os.Environ(), "ssh-keygen", "-f", privatekeypath, "-t", "rsa", "-C", "simulator-key", "-N", "''")
	if *stderr != "" {
		fmt.Println(*stderr)
	}

	return out, err
}

// EnsureKey ensures there is a well-known simulator key available and returns true if it generates a new one or an
// error if any
func EnsureKey() (bool, error) {
	abspath, err := util.ExpandTilde(PrivateKeyPath)
	if err != nil {
		return false, errors.Wrap(err, "Error resolving key path")
	}

	exists, err := util.FileExists(*abspath)
	if err != nil {
		return false, errors.Wrap(err, "Error checking if key already exists")
	}

	// key already exists return that we didn't generate a new one
	if exists {
		return false, nil
	}

	_, err = GenerateKey(*abspath)
	if err != nil {
		return true, errors.Wrap(err, "Error generating key")
	}

	return true, nil
}

// PrivateKeyFile reads the private key at the path supplied and returns the ssh.AuthMethod to use or an error if any
// occurred
func PrivateKeyFile() (ssh.AuthMethod, error) {
	abspath, err := util.ExpandTilde(PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading private key")
	}

	buffer, err := util.Slurp(*abspath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey([]byte(*buffer))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
