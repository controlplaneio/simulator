package simulator_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TfVars_String(t *testing.T) {
	t.Parallel()
	tfv := simulator.NewTfVars("ssh-rsa", "10.0.0.1/16", "test-bucket",
		"latest", "controlplane/simulator-attack", "10.0.0.1/16")
	expected := `access_key = "ssh-rsa"
access_cidr = ["10.0.0.1/16", "10.0.0.1/16"]
attack_container_tag = "latest"
attack_container_repo = "controlplane/simulator-attack"
state_bucket_name = "test-bucket"
`
	assert.Equal(t, expected, tfv.String())
}

func Test_Ensure_TfVarsFile_with_settings(t *testing.T) {
	workDir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	defer os.Remove(workDir)
	require.NoError(t, os.Mkdir(filepath.Join(workDir, "settings"), 0700))

	bastionVarsFile := filepath.Join(workDir, "settings", "bastion.tfvars")
	err = ioutil.WriteFile(bastionVarsFile, []byte("any=content"), 0644)
	require.NoError(t, err)

	err = simulator.EnsureLatestTfVarsFile(workDir, "ssh-rsa", "10.0.0.1/16",
		"test-bucket", "latest", "controlplane/simulator-attack", "10.0.0.1/16, 10.0.0.1/32")
	require.NoError(t, err)
	expected := `access_key = "ssh-rsa"
access_cidr = ["10.0.0.1/16", "10.0.0.1/16", "10.0.0.1/32"]
attack_container_tag = "latest"
attack_container_repo = "controlplane/simulator-attack"
state_bucket_name = "test-bucket"
`
	assert.Equal(t, expected, util.MustSlurp(bastionVarsFile))
}
