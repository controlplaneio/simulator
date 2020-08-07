package simulator

import (
	"strings"

	"github.com/kubernetes-simulator/simulator/pkg/util"
)

// TfVars struct representing the input variables for terraform to create the
// infrastructure
type TfVars struct {
	PublicKey       string
	AccessCIDR      string
	AccessUsername  string
	BucketName      string
	AttackTag       string
	AttackRepo      string
	ExtraCIDRs      string
	GithubUsernames string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, accessUsername, bucketName, attackTag, attackRepo, extraCIDRs, githubUsernames string) TfVars {
	return TfVars{
		PublicKey:       publicKey,
		AccessCIDR:      accessCIDR,
		AccessUsername:  accessUsername,
		BucketName:      bucketName,
		AttackTag:       attackTag,
		AttackRepo:      attackRepo,
		ExtraCIDRs:      extraCIDRs,
		GithubUsernames: githubUsernames,
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

	if tfv.GithubUsernames != "" {
		splitUsernames := strings.Split(tfv.GithubUsernames, ",")
		for i := range splitUsernames {
			splitUsernames[i] = strings.TrimSpace(splitUsernames[i])
		}
		templatedUsernames := strings.Join(splitUsernames, "\", \"")
		tfv.AccessUsername = tfv.AccessUsername + "\", \"" + templatedUsernames
	}

	return "access_key = \"" + tfv.PublicKey + "\"\n" +
		"access_cidr = [\"" + tfv.AccessCIDR + "\"]\n" +
		"access_github_usernames = [\"" + tfv.AccessUsername + "\"]\n" +
		"attack_container_tag = \"" + tfv.AttackTag + "\"\n" +
		"attack_container_repo = \"" + tfv.AttackRepo + "\"\n" +
		"state_bucket_name = \"" + tfv.BucketName + "\"\n"
}

// EnsureLatestTfVarsFile always writes an tfvars file
func EnsureLatestTfVarsFile(tfVarsDir, publicKey, accessCIDR, accessUsername, bucket, attackTag, attackRepo, extraCIDRs, githubUsernames string) error {
	filename := tfVarsDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, accessUsername, bucket, attackTag, attackRepo, extraCIDRs, githubUsernames)

	return util.OverwriteFile(filename, tfv.String())
}
