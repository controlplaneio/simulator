package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
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

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" + "access_cidr = \"" + tfv.AccessCIDR + "\"\n" + "attack_container_tag = \"" + tfv.AttackTag + "\"\n" + "state_bucket_name = \"" + tfv.BucketName + "\"\n"

}

// EnsureLatestTfVarsFile writes an tfvars file if one hasnt already been made
func EnsureLatestTfVarsFile(tfDir, publicKey, accessCIDR, bucket, attackTag string) error {
	filename := tfDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucket, attackTag)

	return util.OverwriteFile(filename, tfv.String())
}
