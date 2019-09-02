package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetAuthMethods(t *testing.T) {
	t.Parallel()
	auths, err := ssh.GetAuthMethods()

	assert.Nil(t, err)
	assert.NotNil(t, auths)
	assert.Equal(t, len(auths), 1)
}
