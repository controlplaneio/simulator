package ssh_test

import (
	"testing"

	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureKey_and_GetAuthMethods(t *testing.T) {
	_, err := ssh.EnsureKey(logrus.New())
	assert.Nil(t, err, "Expected no error ensuring keypair")
	auths, err := ssh.GetAuthMethods()

	assert.Nil(t, err)
	assert.NotNil(t, auths)
	assert.Equal(t, len(auths), 1)
}
