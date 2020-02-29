package simulator

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToArguments_And_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		bastion  net.IP
		master   net.IP
		slaves   []net.IP
		force    bool
		expexted string
	}{
		{
			name:     "base",
			bastion:  net.IPv4(127, 0, 0, 1),
			master:   net.IPv4(127, 0, 0, 1),
			slaves:   []net.IP{net.IPv4(8, 8, 8, 8), net.IPv4(127, 0, 0, 2)},
			expexted: "--master 127.0.0.1 --bastion 127.0.0.1 --nodes 8.8.8.8,127.0.0.2",
		},
		{
			name:     "single slave",
			bastion:  net.IPv4(127, 0, 0, 1),
			master:   net.IPv4(127, 0, 0, 1),
			slaves:   []net.IP{net.IPv4(8, 8, 8, 8)},
			expexted: "--master 127.0.0.1 --bastion 127.0.0.1 --nodes 8.8.8.8",
		},
		{
			name:     "with force",
			bastion:  net.IPv4(127, 0, 0, 1),
			master:   net.IPv4(127, 0, 0, 1),
			slaves:   []net.IP{net.IPv4(8, 8, 8, 8), net.IPv4(127, 0, 0, 2)},
			force:    true,
			expexted: "--master 127.0.0.1 --bastion 127.0.0.1 --nodes 8.8.8.8,127.0.0.2 --force",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			po := PerturbOptions{
				bastion: tt.bastion,
				master:  tt.master,
				slaves:  tt.slaves,
				Force:   tt.force,
			}
			assert.Equal(t, tt.expexted, po.String())
		})
	}
}

func Test_MakePerturbOptions(t *testing.T) {
	t.Parallel()
	tfo := TerraformOutput{
		BastionPublicIP: StringOutput{
			Value: "127.0.0.1",
		},
		MasterNodesPrivateIP: StringSliceOutput{
			Value: []string{"127.0.0.1"},
		},
		ClusterNodesPrivateIP: StringSliceOutput{
			Value: []string{"127.0.0.2"},
		},
	}

	path := "./scenario/test"

	po := PerturbOptions{}
	po.MakePerturbOptions(tfo, path)

	assert.Equal(t, po.bastion.String(), tfo.BastionPublicIP.Value)
	assert.Equal(t, po.master.String(), tfo.MasterNodesPrivateIP.Value[0])
	assert.Equal(t, po.slaves[0].String(), tfo.ClusterNodesPrivateIP.Value[0])

}
