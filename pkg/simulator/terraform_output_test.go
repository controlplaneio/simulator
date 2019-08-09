package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func fixture(name string) string {
	return "../../test/fixtures/" + name
}

func readFixture(name string) string {
	file, err := os.Open(fixture(name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func Test_IsUsable(t *testing.T) {
	tfo := simulator.TerraformOutput{}
	assert.False(t, tfo.IsUsable(), "Empty TerraformOutput was usable")
	tfo.BastionPublicIP.Value = "127.0.0.1"
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only bastion was usable")
	tfo.MasterNodesPrivateIP.Value = []string{"127.0.0.1"}
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only master IP was usable")
	tfo.ClusterNodesPrivateIP.Value = []string{"127.0.0.1"}
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only 1 cluster node IPs was usable")

	tfo.ClusterNodesPrivateIP.Value = []string{"127.0.0.1", "127.0.0.2"}
	assert.True(t, tfo.IsUsable(), "Complete TerraformOutput was not usable")
}

func Test_ToSSHConfig(t *testing.T) {
	tfo := simulator.TerraformOutput{
		BastionPublicIP: simulator.StringOutput{
			Sensitive: false,
			Type:      "string",
			Value:     "8.8.8.8",
		},
		MasterNodesPrivateIP: simulator.StringSliceOutput{
			Sensitive: false,
			Type:      []interface{}{},
			Value:     []string{"127.0.0.1"},
		},
		ClusterNodesPrivateIP: simulator.StringSliceOutput{
			Sensitive: false,
			Type:      []interface{}{},
			Value:     []string{"127.0.0.2", "127.0.0.3"},
		},
	}
	expected := `Host bastion 8.8.8.8
  Hostname 8.8.8.8
  User root
  RequestTTY force
  IdentityFile ~/.ssh/cp_simulator_rsa
  UserKnownHostsFile ~/.ssh/cp_simulator_known_hosts
Host master-0 127.0.0.1
  Hostname 127.0.0.1
  User root
  RequestTTY force
  IdentityFile ~/.ssh/cp_simulator_rsa
  UserKnownHostsFile ~/.ssh/cp_simulator_known_hosts
  ProxyJump bastion
Host node-0 127.0.0.2
  Hostname 127.0.0.2
  User root
  RequestTTY force
  IdentityFile ~/.ssh/cp_simulator_rsa
  UserKnownHostsFile ~/.ssh/cp_simulator_known_hosts
  ProxyJump bastion
Host node-1 127.0.0.3
  Hostname 127.0.0.3
  User root
  RequestTTY force
  IdentityFile ~/.ssh/cp_simulator_rsa
  UserKnownHostsFile ~/.ssh/cp_simulator_known_hosts
  ProxyJump bastion
`

	out, err := tfo.ToSSHConfig()
	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "Got nil output")

	assert.Equal(t, *out, expected, "SSH config was not correct")
}

func Test_ParseTerraformOutput(t *testing.T) {
	t.Parallel()
	output := readFixture("valid-tf-output.json")

	tfOutput, err := simulator.ParseTerraformOutput(output)

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfOutput, "Output was nil")
	assert.Equal(t, tfOutput.BastionPublicIP.Value, "34.244.109.234", "Bastion IP was wrong")
	assert.Equal(t, len(tfOutput.ClusterNodesPrivateIP.Value), 1, "Didn't get 1 node IP")
	assert.Equal(t, tfOutput.ClusterNodesPrivateIP.Value[0], "172.31.2.19")
	assert.Equal(t, len(tfOutput.MasterNodesPrivateIP.Value), 1, "Didn't get 1 master IP")
	assert.Equal(t, tfOutput.MasterNodesPrivateIP.Value[0], "172.31.2.167")
}
