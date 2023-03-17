package simulator

import (
	"fmt"
	"strings"

	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pelletier/go-toml"
)

// TfVars struct representing the input variables for terraform to create the
// infrastructure
type TfVars struct {
	PublicKey       string   `toml:"access_key"`
	AccessCIDR      []string `toml:"access_cidr"`
	BucketName      string   `toml:"state_bucket_name"`
	AttackTag       string   `toml:"attack_container_tag"`
	AttackRepo      string   `toml:"attack_container_repo"`
	ExtraCIDRs      []string `toml:"-"`
	GithubUsernames []string `toml:"access_github_usernames"`
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey, accessCIDR, bucketName, attackTag, attackRepo, extraCIDRs, githubUsernames string) TfVars {
	var splitCIDRs []string
	if extraCIDRs != "" {
		splitCIDRs = strings.Split(extraCIDRs, ",")
		for i := range splitCIDRs {
			splitCIDRs[i] = strings.TrimSpace(splitCIDRs[i])
		}
	}

	var allCIDRs []string
	if accessCIDR != "" {
		allCIDRs = append([]string{accessCIDR}, splitCIDRs...)
	}

	var splitUsernames []string
	if githubUsernames != "" {
		splitUsernames = strings.Split(githubUsernames, ",")
		for i := range splitUsernames {
			splitUsernames[i] = strings.TrimSpace(splitUsernames[i])
		}
	}

	return TfVars{
		PublicKey:       publicKey,
		AccessCIDR:      allCIDRs,
		BucketName:      bucketName,
		AttackTag:       attackTag,
		AttackRepo:      attackRepo,
		ExtraCIDRs:      splitCIDRs,
		GithubUsernames: splitUsernames,
	}
}

func (tfv *TfVars) Marshal() ([]byte, error) {
	cfgBytes, err := toml.Marshal(tfv)
	if err != nil {
		return nil, fmt.Errorf("error marshaling TFVars toml")
	}

	return cfgBytes, nil
}

// EnsureLatestTfVarsFile always writes an tfvars file
func EnsureLatestTfVarsFile(tfVarsDir, publicKey, accessCIDR, bucket, attackTag, attackRepo, extraCIDRs, githubUsernames string) error {
	filename := tfVarsDir + "/settings/bastion.tfvars"
	tfv := NewTfVars(publicKey, accessCIDR, bucket, attackTag, attackRepo, extraCIDRs, githubUsernames)

	contents, err := tfv.Marshal()
	if err != nil {
		return err
	}

	return util.OverwriteFile(filename, string(contents))
}
