package simulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"text/template"
)

// StringOutput is a struct representing an output from terraform that contains
// a string
type StringOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

// StringSliceOutput is a struct representing an output from terraform that
// contains a slice of strings
type StringSliceOutput struct {
	Sensitive bool          `json:"sensitive"`
	Type      []interface{} `json:"type"`
	Value     []string      `json:"value"`
}

// TerraformOutput is a struct representing the expected output variables from
// the terraform script
type TerraformOutput struct {
	BastionPublicIP       StringOutput      `json:"bastion_public_ip"`
	ClusterNodesPrivateIP StringSliceOutput `json:"cluster_nodes_private_ip"`
	MasterNodesPrivateIP  StringSliceOutput `json:"master_nodes_private_ip"`
}

var bastionConfigTmplSrc = `Host bastion {{.Hostname}}
  Hostname {{.Hostname}}
  User root
  RequestTTY force
  IdentityFile {{.KeyFilePath}}
  UserKnownHostsFile {{.KnownHostsFilePath}}
`
var k8sConfigTmplSrc = `Host {{.Alias}} {{.Hostname}}
  Hostname {{.Hostname}}
  User root
  RequestTTY force
  IdentityFile {{.KeyFilePath}}
  UserKnownHostsFile {{.KnownHostsFilePath}}
  ProxyJump bastion
`

// SSHConfig represents the values needed to produce a config block to allow
// SSH to the private kubernetes nodes via the bastion
type SSHConfig struct {
	Alias              string
	Hostname           string
	KeyFilePath        string
	KnownHostsFilePath string
	BastionIP          string
}

// ToSSHConfig produces the SSH config
func (tfo *TerraformOutput) ToSSHConfig() (*string, error) {
	bastionConfigTmpl, err := template.New("bastion-ssh-config").Parse(bastionConfigTmplSrc)
	k8sConfigTmpl, err := template.New("k8s-ssh-config").Parse(k8sConfigTmplSrc)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ssh config template")
	}

	var buf bytes.Buffer
	bastionCfg := SSHConfig{
		Alias:              "bastion",
		Hostname:           tfo.BastionPublicIP.Value,
		KeyFilePath:        "~/.ssh/cp_simulator_rsa",
		KnownHostsFilePath: "~/.ssh/cp_simulator_known_hosts",
	}
	err = bastionConfigTmpl.Execute(&buf, bastionCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error populating ssh bastion config template with %+v", bastionCfg)
	}

	for i, ip := range tfo.MasterNodesPrivateIP.Value {
		c := SSHConfig{
			Alias:              fmt.Sprintf("master-%d", i),
			Hostname:           ip,
			KeyFilePath:        "~/.ssh/cp_simulator_rsa",
			KnownHostsFilePath: "~/.ssh/cp_simulator_known_hosts",
			BastionIP:          tfo.BastionPublicIP.Value,
		}

		err = k8sConfigTmpl.Execute(&buf, c)
		if err != nil {
			return nil, errors.Wrapf(err, "Error populating ssh master config template with %+v", c)
		}
	}

	for i, ip := range tfo.ClusterNodesPrivateIP.Value {
		c := SSHConfig{
			Alias:              fmt.Sprintf("node-%d", i),
			Hostname:           ip,
			KeyFilePath:        "~/.ssh/cp_simulator_rsa",
			KnownHostsFilePath: "~/.ssh/cp_simulator_known_hosts",
			BastionIP:          tfo.BastionPublicIP.Value,
		}
		err = k8sConfigTmpl.Execute(&buf, c)
		if err != nil {
			return nil, errors.Wrapf(err, "Error populating ssh node config template with %+v", c)
		}
	}

	var output = string(buf.Bytes())
	return &output, nil
}

// IsUsable checks whether the TerraformOutput has all the necessary
// information to be converted for use with perturb
func (tfo *TerraformOutput) IsUsable() bool {
	return tfo.BastionPublicIP.Value != "" && len(tfo.MasterNodesPrivateIP.Value) == 1 && len(tfo.ClusterNodesPrivateIP.Value) == 2
}

// ParseTerraformOutput takes a string containing the stdout from `terraform
// output -json` and returns a TerraformOutput struct
func ParseTerraformOutput(output string) (*TerraformOutput, error) {
	out := TerraformOutput{}
	err := json.Unmarshal([]byte(output), &out)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshalling json %s", output)
	}

	return &out, nil
}
