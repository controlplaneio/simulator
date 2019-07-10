package util

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

// Base64PrivateKey returns a pointer to a string containing the base64 encoded private key or an error
func Base64PrivateKey(name string) (*string, error) {
	keypath, err := ExpandTilde(name)
	if err != nil {
		return nil, err
	}

	keymaterial, err := Slurp(*keypath)
	if err != nil {
		return nil, err
	}

	encodedkey := base64.StdEncoding.EncodeToString([]byte(*keymaterial))

	return &encodedkey, nil
}

// PublicKey reads the public key and return a pointer to a string with its contents or any error
func PublicKey() (*string, error) {
	publicKeyPath, err := ExpandTilde(PublicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", PublicKeyPath)
	}

	publickey, err := Slurp(*publicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", PublicKeyPath)
	}

	ret := strings.Trim(*publickey, "\n")
	return &ret, nil
}

// KeyScan runs ssh-keyscan silently against the provided bastion address. It returns a pointer to a string containing
// its buffered stdout or an error if any occurred
func KeyScan(bastion string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := RunSilently(wd, os.Environ(), "ssh-keyscan", "-H", bastion)
	return out, err
}

// GenerateKey runs ssh-keygen silently to create an SSH key with the same provided using preconfigured settings
// It returns a pointer to a string containing the buffered stdout or an error if any occurred
func GenerateKey(privatekeypath string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, stderr, err := RunSilently(wd, os.Environ(), "ssh-keygen", "-f", privatekeypath, "-t", "rsa", "-C", "simulator-key", "-N", "''")
	if *stderr != "" {
		fmt.Println(*stderr)
	}

	return out, err
}

// PrivateKeyPath is the path to the key to be generated and used by simulator
const PrivateKeyPath = "~/.ssh/cp_simulator_rsa"

// PublicKeyPath is the path to the key to be generated and used by simulator
const PublicKeyPath = PrivateKeyPath + ".pub"

// EnsureKey ensures there is a well-known simulator key available and returns true if it generates a new one or an
// error if any
func EnsureKey() (bool, error) {
	abspath, err := ExpandTilde(PrivateKeyPath)
	if err != nil {
		return false, errors.Wrap(err, "Error resolving key path")
	}

	exists, err := FileExists(*abspath)
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
func PrivateKeyFile(file string) (ssh.AuthMethod, error) {
	abspath, err := ExpandTilde(keypath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s when falling back to key", file)
	}

	buffer, err := Slurp(*abspath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey([]byte(*buffer))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
