package ssh_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureKey_and_GetAuthMethods(t *testing.T) {
	os.Remove(ssh.PublicKeyPath)
	ls := ssh.LocalStateProvider{}
	kp, err := ls.GetSSHKeyPair()
	fmt.Printf("%-v", kp)
	assert.Nil(t, err, "Expected no error ensuring keypair")
	auths, err := ssh.GetAuthMethods(*kp)

	assert.Nil(t, err)
	assert.NotNil(t, auths)
	assert.Equal(t, len(auths), 1)
}
