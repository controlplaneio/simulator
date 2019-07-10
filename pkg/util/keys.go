package util

import (
	"encoding/base64"
	"github.com/pkg/errors"
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
