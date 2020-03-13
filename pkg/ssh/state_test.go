package ssh_test

import (
	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetSSHKeyPair(t *testing.T) {
	ls := ssh.LocalStateProvider{}

	kp, err := ls.GetSSHKeyPair()
	assert.Nil(t, err, "Expected no error ensuring keypair")
	assert.NotNil(t, kp, "Expected a key pair")
}
