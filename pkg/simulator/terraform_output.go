package simulator

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"os/user"
	"text/template"
)

// StringOutput is a struct representing an output from terraform that contains a string
type StringOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

// StringSliceOutput is a struct representing an output from terraform that contains a slice of strings
type StringSliceOutput struct {
	Sensitive bool          `json:"sensitive"`
	Type      []interface{} `json:"type"`
	Value     []string      `json:"value"`
}

// TerraformOutput is a struct representing the expected output variables from the terraform script
type TerraformOutput struct {
	BastionPublicIP       StringOutput      `json:"bastion_public_ip"`
	ClusterNodesPrivateIP StringSliceOutput `json:"cluster_nodes_private_ip"`
	MasterNodesPrivateIP  StringSliceOutput `json:"master_nodes_private_ip"`
}

var sshConfigTmplSrc = `Host {{.Hostname}}
  IdentityFile {{.KeyFilePath}}
  ProxyCommand ssh root@{{.BastionIP}} -W %h:%p
`

// SSHConfig represents the values needed to produce a config block to allow
// SSH to the private kubernetes nodes via the bastion
type SSHConfig struct {
	Hostname    string
	KeyFilePath string
	User        string
	BastionIP   string
}

// ToSSHConfig produces the SSH config
func (tfo *TerraformOutput) ToSSHConfig() (*string, error) {
	var sshConfigTmpl, err = template.New("ssh-config").Parse(sshConfigTmplSrc)
	if err != nil {
		return nil, err
	}

	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get current user for generating sshconfig")
	}

	var buf bytes.Buffer
	for _, ip := range tfo.MasterNodesPrivateIP.Value {
		c := SSHConfig{
			Hostname:    ip,
			KeyFilePath: "~/.ssh/id_rsa.pub",
			User:        u.Username,
			BastionIP:   tfo.BastionPublicIP.Value,
		}

		err = sshConfigTmpl.Execute(&buf, c)
		if err != nil {
			return nil, errors.Wrapf(err, "Error populating ssh config template with %+v", c)
		}
	}

	for _, ip := range tfo.ClusterNodesPrivateIP.Value {
		c := SSHConfig{
			Hostname:    ip,
			KeyFilePath: "~/.ssh/id_rsa.pub",
			User:        u.Username,
			BastionIP:   tfo.BastionPublicIP.Value,
		}

		err = sshConfigTmpl.Execute(&buf, c)
		if err != nil {
			return nil, errors.Wrapf(err, "Error populating ssh config template with %+v", c)
		}
	}

	var output = string(buf.Bytes())
	return &output, nil
}

// IsUsable checks whether the TerraformOutput has all the necessary information to be converted for use with perturb
func (tfo *TerraformOutput) IsUsable() bool {
	return tfo.BastionPublicIP.Value != "" && len(tfo.MasterNodesPrivateIP.Value) == 1 && len(tfo.ClusterNodesPrivateIP.Value) == 2
}

// ParseTerraformOutput takes a string containing the stdout from `terraform output -json` and returns a TerraformOutput
// struct
func ParseTerraformOutput(output string) (*TerraformOutput, error) {
	out := TerraformOutput{}
	err := json.Unmarshal([]byte(output), &out)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshalling json %s", output)
	}

	return &out, nil
}
