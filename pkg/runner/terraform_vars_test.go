package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/runner"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TfVars_String(t *testing.T) {
	t.Parallel()
	tfv := runner.NewTfVars("ssh-rsa", "10.0.0.1/16")
	expected := `access_key = "ssh-rsa"
access_cidr = "10.0.0.1/16"
`
	assert.Equal(t, tfv.String(), expected)
}

func Test_Ensure_TfVarsFile_no_settings(t *testing.T) {
	tfDir := fixture("noop-tf-dir")

	err := runner.EnsureTfVarsFile(tfDir, "ssh-rsa", "10.0.0.1/16")
	assert.Nil(t, err, "Got an error")

	exists, err := util.FileExists(tfDir + "/settings/bastion.tfVars")
	assert.Nil(t, err, "Got an error checking file had been written")
	assert.True(t, exists, "File wasn't created")
}

func Test_Ensure_TfVarsFile_with_settings(t *testing.T) {
	tfDir := fixture("tf-dir-with-settings")
	varsFile := tfDir + "/settings/bastion.tfVars"

	err := runner.EnsureTfVarsFile(tfDir, "ssh-rsa", "10.0.0.1/16")
	assert.Nil(t, err, "Got an error")

	assert.Equal(t, util.MustSlurp(varsFile), "test = true\n")
}
