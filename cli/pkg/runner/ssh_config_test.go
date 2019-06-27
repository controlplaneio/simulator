package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateSshConfig(t *testing.T) {
	tfo := runner.TerraformOutput{
		BastionPublicIP: runner.StringOutput{
			Sensitive: false,
			Type:      "string",
			Value:     "8.8.8.8",
		},
		MasterNodesPrivateIP: runner.StringSliceOutput{
			Sensitive: false,
			Type:      []interface{}{},
			Value:     []string{"127.0.0.1"},
		},
	}
	const expected = `Host 127.0.0.1
  IdentityFile ~/.ssh/id_rsa.pub
  ProxyCommand ssh raoul@8.8.8.8 -W %h:%p
`

	out, err := runner.CreateSshConfig(tfo)
	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "Got nil output")

	assert.Equal(t, *out, expected, "SSH config was not correct")
}
