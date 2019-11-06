package simulator

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TfVars struct representing the input variables for terraform to create the
// infrastructure
type TfVars struct {
	PublicKey  string
	AccessCIDR string
	BucketName string
	AttackTag  string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, bucketName, attackTag string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
		BucketName: bucketName,
		AttackTag:  attackTag,
	}
}

func writeProvidersFile(tfDir, bucket string) error {
	providerspath := filepath.Join(tfDir, "providers.tf")
	input, err := ioutil.ReadFile(providerspath)
	if err != nil {
		return errors.Wrapf(err, "Error reading providers file %s", providerspath)
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "bucket = ") {
			lines[i] = fmt.Sprintf("    bucket = \"%s\"", bucket)
		}
	}
	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(providerspath, []byte(output), 0644)
	if err != nil {
		return errors.Wrapf(err, "Error writing providers file %s", providerspath)
	}

	return nil
}

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" + "access_cidr = \"" + tfv.AccessCIDR + "\"\n" + "attack_container_tag = \"" + tfv.AttackTag + "\"\n"

}

// EnsureLatestTfVarsFile writes an tfvars file if one hasnt already been made
func EnsureLatestTfVarsFile(tfDir, publicKey, accessCIDR, bucket, attackTag string) error {
	filename := tfDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucket, attackTag)

	err := writeProvidersFile(tfDir, bucket)
	if err != nil {
		return errors.Wrap(err, "Error saving bucket name")

	}

	return util.OverwriteFile(filename, tfv.String())
}
