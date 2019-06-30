package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/runner"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func Test_ToArguments_And_String(t *testing.T) {
	t.Parallel()
	po := runner.PerturbOptions{
		Master: net.IPv4(127, 0, 0, 1),
		Slaves: []net.IP{net.IPv4(8, 8, 8, 8), net.IPv4(127, 0, 0, 2)},
	}

	assert.Equal(t, po.String(), "--master 127.0.0.1 --slaves 8.8.8.8,127.0.0.2")
}

func Test_MakePerturbOptions(t *testing.T) {
	t.Parallel()
	tfo := runner.TerraformOutput{
		MasterNodesPrivateIP: runner.StringSliceOutput{
			Value: []string{"127.0.0.1"},
		},
		ClusterNodesPrivateIP: runner.StringSliceOutput{
			Value: []string{"127.0.0.2"},
		},
	}

	path := "./scenario/test"

	po := runner.MakePerturbOptions(tfo, path)

	assert.Equal(t, po.Master.String(), tfo.MasterNodesPrivateIP.Value[0])
	assert.Equal(t, po.Slaves[0].String(), tfo.ClusterNodesPrivateIP.Value[0])

}

func Test_Perturb(t *testing.T) {
	os.Setenv("SIMULATOR_MANIFEST_PATH", fixture("noop-perturb"))
	po := runner.PerturbOptions{}
	_, err := runner.Perturb(&po)

	assert.Nil(t, err, "Got an error")
}
