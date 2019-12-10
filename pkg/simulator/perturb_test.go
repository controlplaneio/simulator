package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func Test_ToArguments_And_String(t *testing.T) {
	t.Parallel()
	po := simulator.PerturbOptions{
		Bastion: net.IPv4(127, 0, 0, 1),
		Master:  net.IPv4(127, 0, 0, 1),
		Slaves:  []net.IP{net.IPv4(8, 8, 8, 8), net.IPv4(127, 0, 0, 2)},
	}

	assert.Equal(t, po.String(), "--master 127.0.0.1 --bastion 127.0.0.1 --nodes 8.8.8.8,127.0.0.2")
}

func Test_MakePerturbOptions(t *testing.T) {
	t.Parallel()
	tfo := simulator.TerraformOutput{
		BastionPublicIP: simulator.StringOutput{
			Value: "127.0.0.1",
		},
		MasterNodesPrivateIP: simulator.StringSliceOutput{
			Value: []string{"127.0.0.1"},
		},
		ClusterNodesPrivateIP: simulator.StringSliceOutput{
			Value: []string{"127.0.0.2"},
		},
	}

	path := "./scenario/test"

	po := simulator.MakePerturbOptions(tfo, path)

	assert.Equal(t, po.Bastion.String(), tfo.BastionPublicIP.Value)
	assert.Equal(t, po.Master.String(), tfo.MasterNodesPrivateIP.Value[0])
	assert.Equal(t, po.Slaves[0].String(), tfo.ClusterNodesPrivateIP.Value[0])

}

func Test_Perturb(t *testing.T) {
	os.Setenv("SIMULATOR_SCENARIOS_DIR", fixture("noop-perturb"))
	po := simulator.PerturbOptions{}
	_, err := simulator.Perturb(&po)

	assert.Nil(t, err, "Got an error")
}
