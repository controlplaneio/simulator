package ssh_test

import (
	"testing"

	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureKey_and_GetAuthMethods(t *testing.T) {
	ls := ssh.LocalState{}
	kp, err := ls.GetSSHKeyPair()
	assert.Nil(t, err, "Expected no error ensuring keypair")
	auths, err := ssh.GetAuthMethods(*kp)

	assert.Nil(t, err)
	assert.NotNil(t, auths)
	assert.Equal(t, len(auths), 1)
}
