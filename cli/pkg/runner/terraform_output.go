package runner

import (
	"encoding/json"
	"github.com/pkg/errors"
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
