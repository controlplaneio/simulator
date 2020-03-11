package simulator

import (
	"strings"

	"github.com/kubernetes-simulator/simulator/pkg/util"
)

// TfVars struct representing the input variables for terraform to create the
// infrastructure
type TfVars struct {
	PublicKey  string
	AccessCIDR string
	BucketName string
	AttackTag  string
	AttackRepo string
	ExtraCIDRs string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, bucketName, attackTag, attackRepo, extraCIDRs string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
		BucketName: bucketName,
		AttackTag:  attackTag,
		AttackRepo: attackRepo,
		ExtraCIDRs: extraCIDRs,
	}
}

func (tfv *TfVars) String() string {
	if tfv.ExtraCIDRs != "" {
		splitCIDRs := strings.Split(tfv.ExtraCIDRs, ",")
		for i := range splitCIDRs {
			splitCIDRs[i] = strings.TrimSpace(splitCIDRs[i])
		}
		templatedCIDRs := strings.Join(splitCIDRs, "\", \"")
		tfv.AccessCIDR = tfv.AccessCIDR + "\", \"" + templatedCIDRs
	}

	return "access_key = \"" + tfv.PublicKey + "\"\n" +
		"access_cidr = [\"" + tfv.AccessCIDR + "\"]\n" +
		"attack_container_tag = \"" + tfv.AttackTag + "\"\n" +
		"attack_container_repo = \"" + tfv.AttackRepo + "\"\n" +
		"state_bucket_name = \"" + tfv.BucketName + "\"\n"
}

// EnsureLatestTfVarsFile always writes an tfvars file
func EnsureLatestTfVarsFile(tfVarsDir, publicKey, accessCIDR, bucket, attackTag, attackRepo, extraCIDRs string) error {
	filename := tfVarsDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucket, attackTag, attackRepo, extraCIDRs)

	return util.OverwriteFile(filename, tfv.String())
}
