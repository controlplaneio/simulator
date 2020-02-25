package ssh_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EnsureKnownHosts(t *testing.T) {
	t.Skip("Need to orchestrate an SSH server for testing")
	err := ssh.EnsureKnownHosts("localhost", logrus.New())
	assert.NotNil(t, err, "Expected no error")

}
