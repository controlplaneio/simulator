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
	AttackRepo string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, bucketName, attackTag, attackRepo string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
		BucketName: bucketName,
		AttackTag:  attackTag,
		AttackRepo: attackRepo,
	}
}

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" +
		"access_cidr = \"" + tfv.AccessCIDR + "\"\n" +
		"attack_container_tag = \"" + tfv.AttackTag + "\"\n" +
		"attack_container_repo = \"" + tfv.AttackRepo + "\"\n" +
		"state_bucket_name = \"" + tfv.BucketName + "\"\n"

}

// EnsureLatestTfVarsFile always writes an tfvars file
func EnsureLatestTfVarsFile(tfVarsDir, publicKey, accessCIDR, bucket, attackTag, attackRepo string) error {
	filename := tfVarsDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucket, attackTag, attackRepo)

	return util.OverwriteFile(filename, tfv.String())
}
