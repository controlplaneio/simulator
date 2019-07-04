package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
)

// TfVars struct representing the input variables for terraform to create the infrastructure
type TfVars struct {
	PublicKey  string
	AccessCIDR string
	BucketName string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, bucketName string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
		BucketName: "simulator-" + bucketName,
	}
}

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" + "access_cidr = \"" + tfv.AccessCIDR + "\"\n" + "s3_bucket_name = \"" + tfv.BucketName + "\"\n"
}

// EnsureTfVarsFile writes an tfvars file if one hasnt already been made
func EnsureTfVarsFile(tfDir, publicKey, accessCIDR, bucketName string) error {
	filename := tfDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucketName)

	_, err := util.EnsureFile(filename, tfv.String())
	return err
}
