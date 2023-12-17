package acceptance_test

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"github.com/controlplaneio/simulator/v2/core/aws"
	"github.com/controlplaneio/simulator/v2/utils/ssh"
)

const (
	bastion = "bastion"
)

// TestSimulatorWorkspaceCreatesAccessibleCluster Performs the following steps to acceptance test the Terraform
// 1. Terraform the infrastructure
// 2. Using the AdminBundle, ssh into the Bastion and execute 'kubectl get nodes'
// 3. Verify the master and two nodes are initialised, but not ready (has no cluster network yet)
func TestSimulatorWorkspaceCreatesAccessibleCluster(t *testing.T) {
	if os.Getenv("RUN_ACCEPTANCE_TEST") != strings.ToLower("yes") {
		t.Skip("Skipping acceptance tests as RUN_ACCEPTANCE_TEST is not set to 'yes'")
	}

	t.Parallel()

	currentDir, err := os.Getwd()
	assert.NoError(t, err)

	workingDir, err := os.MkdirTemp("", "workspace")
	assert.NoError(t, err)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(workingDir)

	adminBundleDir := filepath.Join(workingDir, "admin")
	playerBundleDir := filepath.Join(workingDir, "player")
	terraformDir := filepath.Join(currentDir, "..", "terraform")
	ansiblePlaybooksDir := filepath.Join(currentDir, "..", "ansible", "playbooks")
	ansibleRolesDir := filepath.Join(currentDir, "..", "ansible", "roles")
	bucket := randomstring.EnglishFrequencyString(18)

	err = copy.Copy(terraformDir, workingDir)
	assert.NoError(t, err)

	ctx := context.Background()

	s3Client, err := aws.NewS3Client(ctx)
	assert.NoError(t, err)

	err = s3Client.Create(ctx, bucket)
	assert.NoError(t, err)
	defer func(s3 *aws.S3Client, ctx context.Context, name string) {
		err = s3.Delete(ctx, name)
		assert.NoError(t, err)
	}(s3Client, ctx, bucket)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: filepath.Join(workingDir, "workspaces", "simulator"),
		Vars: map[string]interface{}{
			"name":                 "simulator-test",
			"ansible_playbook_dir": ansiblePlaybooksDir,
			"ansible_roles_dir":    ansibleRolesDir,
			"admin_bundle_dir":     adminBundleDir,
			"player_bundle_dir":    playerBundleDir,
		},
		BackendConfig: map[string]interface{}{
			"bucket": bucket,
			"key":    "terraform.tfstate",
		},
		Reconfigure: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	_, err = terraform.InitAndApplyE(t, terraformOptions)
	assert.NoError(t, err)

	sshClient, err := ssh.NewClient(
		adminBundleDir,
		"simulator_config",
		[]string{
			ssh.KeyAlgoED25519,
		},
		15*time.Second,
	)
	assert.NoError(t, err)

	output, err := sshClient.Execute("kubectl get nodes")
	assert.NoError(t, err)

	assert.Regexp(t, regexp.MustCompile("master-1.*NotReady"), string(output))
	assert.Regexp(t, regexp.MustCompile("node-1.*NotReady"), string(output))
	assert.Regexp(t, regexp.MustCompile("node-2.*NotReady"), string(output))
}
