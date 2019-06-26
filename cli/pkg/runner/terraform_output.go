package runner

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type StringOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type StringSliceOutput struct {
	Sensitive bool          `json:"sensitive"`
	Type      []interface{} `json:"type"`
	Value     []string      `json:"value"`
}

type TerraformOutput struct {
	BastionPublicIP       StringOutput      `json:"bastion_public_ip"`
	ClusterNodesPrivateIP StringSliceOutput `json:"cluster_nodes_private_ip"`
	MasterNodesPrivateIP  StringSliceOutput `json:"master_nodes_private_ip"`
}

func ParseTerraformOutput(output string) (*TerraformOutput, error) {
	out := TerraformOutput{}
	err := json.Unmarshal([]byte(output), &out)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshalling json %s", output)
	}

	return &out, nil
}
