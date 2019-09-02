package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrivateKeyFile_returns_AuthMethod(t *testing.T) {
	_, err := ssh.EnsureKey()
	assert.Nil(t, err, "Expected no error ensuring keypair")

	auth, err := ssh.PrivateKeyFile()
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, auth)
}

func Test_PublicKey_returns_Key_file_contents(t *testing.T) {
	_, err := ssh.EnsureKey()
	assert.Nil(t, err, "Expected no error ensuring keypair")

	key, err := ssh.PublicKey()
	assert.Nil(t, err, "Expected no error getting public key")
	assert.NotNil(t, key)
}
