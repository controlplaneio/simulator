package ssh

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SSHPrivateKey represents an SSH PrivateKey
type SSHPrivateKey string

// ToBase64 returns the base64 encoded string representation of the private key
func (pk SSHPrivateKey) ToBase64() string {
	fmt.Println("Encoding Private key as base64")
	encoded := base64.StdEncoding.EncodeToString([]byte(string(pk)))
	fmt.Println(encoded)
	return encoded
}

// ToAuthMethod converts the SSHPrivateKey to an AuthMethod required by the
// crypto ssh library
func (pk SSHPrivateKey) ToAuthMethod() (ssh.AuthMethod, error) {
	key, err := ssh.ParsePrivateKey([]byte(string(pk)))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing private key file")
	}
	return ssh.PublicKeys(key), nil
}

// SSHPublicKey represents an SSH PublicKey
type SSHPublicKey string

// KeyPair has an SSH Private and Public key pair
type KeyPair struct {
	PublicKey  SSHPublicKey
	PrivateKey SSHPrivateKey
}
