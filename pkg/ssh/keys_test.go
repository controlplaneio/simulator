package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrivateKeyFile_returns_AuthMethod(t *testing.T) {
	t.Parallel()

	_, err := ssh.EnsureKey()
	assert.Nil(t, err, "Expected no error ensuring keypair")

	auth, err := ssh.PrivateKeyFile()
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, auth)
}
