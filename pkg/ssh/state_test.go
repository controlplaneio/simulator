package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetSSHKeyPair(t *testing.T) {
	ls := ssh.LocalState{}

	kp, err := ls.GetSSHKeyPair()
	assert.Nil(t, err, "Expected no error ensuring keypair")
	assert.NotNil(t, kp, "Expected a key pair")
}

func Test_NewSSHKeyPair(t *testing.T) {
	ls := ssh.LocalState{}

	kp, err := ls.NewSSHKeyPair()
	assert.Nil(t, err, "Expected no error ensuring keypair")
	assert.NotNil(t, kp, "Expected a key pair")
}
