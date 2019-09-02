package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EnsureKey_and_GetAuthMethods(t *testing.T) {
	t.Parallel()
	_, err := ssh.EnsureKey()
	assert.Nil(t, err, "Expected no error ensure keypair")
	auths, err := ssh.GetAuthMethods()

	assert.Nil(t, err)
	assert.NotNil(t, auths)
	assert.Equal(t, len(auths), 1)
}
