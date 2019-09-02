package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EnsureSSHConfig(t *testing.T) {
	t.Parallel()

	err := ssh.EnsureSSHConfig("test")
	assert.Nil(t, err, "Expected no error ensuring SSH config")
	abspath, err := util.ExpandTilde(ssh.SSHConfigPath)
	assert.Nil(t, err, "Expected no error expanding SSH config path")
	config, err := util.Slurp(*abspath)
	assert.Nil(t, err, "Expected no error reading ssh config")
	assert.NotNil(t, config, "Config was nil")
	assert.Equal(t, *config, "test", "Config contents was not correct")
}
