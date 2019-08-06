package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TfVars_String(t *testing.T) {
	t.Parallel()
	tfv := simulator.NewTfVars("ssh-rsa", "10.0.0.1/16", "test")
	expected := `access_key = "ssh-rsa"
access_cidr = "10.0.0.1/16"
`
	assert.Equal(t, tfv.String(), expected)
}

func Test_Ensure_TfVarsFile_with_settings(t *testing.T) {
	tfDir := fixture("tf-dir-with-settings")
	varsFile := tfDir + "/settings/bastion.tfVars"

	err := simulator.EnsureLatestTfVarsFile(tfDir, "ssh-rsa", "10.0.0.1/16", "test")
	assert.Nil(t, err, "Got an error")

	assert.Equal(t, util.MustSlurp(varsFile), "test = true\n")
}
