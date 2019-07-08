package util

import (
	"encoding/base64"
)

// Base64PrivateKey returns a pointer to a string containing the base64 encoded private key or an error
func Base64PrivateKey(name string) (*string, error) {
	keypath, err := ExpandTilde("~/.ssh/" + name)
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
