package util

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"os"
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
func PublicKey(name string) (*string, error) {
	publicKeyPath, err := ExpandTilde(name)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s", name)
	}

	publickey, err := Slurp(*publicKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading %s", name)
	}

	return publickey, nil
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
func GenerateKey(keyname string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := RunSilently(wd, os.Environ(), "ssh-keygen", "-f", keyname, "-t", "rsa", "-C", "''", "-N", "''")
	return out, err
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
