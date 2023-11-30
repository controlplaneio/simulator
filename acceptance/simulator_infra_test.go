package acceptance_test

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/kevinburke/ssh_config"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/controlplaneio/simulator/v2/core/aws"
)

const (
	bastion = "bastion"
)

// TestSimulatorWorkspaceCreatesAccessibleCluster Performs the following steps to acceptance test the Terraform
// 1. Terraform the infrastructure
// 2. Using the AdminBundle, ssh into the Bastion and execute 'kubectl get nodes'
// 3. Verify the master and two nodes are initialised, but not ready (has no cluster network yet)
func TestSimulatorWorkspaceCreatesAccessibleCluster(t *testing.T) {
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

	sshConfigFile, err := os.Open(filepath.Join(adminBundleDir, "simulator_config"))
	assert.NoError(t, err)

	sshConfig, err := ssh_config.Decode(sshConfigFile)
	assert.NoError(t, err)

	user, _ := sshConfig.Get(bastion, "User")
	host, _ := sshConfig.Get(bastion, "Hostname")
	idFile, _ := sshConfig.Get(bastion, "IdentityFile")
	knownHostFile, _ := sshConfig.Get(bastion, "UserKnownHostsFile")

	key, err := os.ReadFile(filepath.Join(adminBundleDir, idFile))
	assert.NoError(t, err)

	signer, err := ssh.ParsePrivateKey(key)
	assert.NoError(t, err)

	hostKeyCallback, err := knownhosts.New(filepath.Join(adminBundleDir, knownHostFile))
	assert.NoError(t, err)

	sshClientConf := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoED25519,
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         15 * time.Second,
	}

	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, "22"), sshClientConf)
	assert.NoError(t, err)
	defer func(conn *ssh.Client) {
		_ = conn.Close()
	}(conn)

	session, err := conn.NewSession()
	assert.NoError(t, err)
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	stdout, err := session.StdoutPipe()
	assert.NoError(t, err)

	err = session.Start("kubectl get nodes")
	assert.NoError(t, err)

	output, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	assert.Regexp(t, regexp.MustCompile("master-1.*NotReady"), string(output))
	assert.Regexp(t, regexp.MustCompile("node-1.*NotReady"), string(output))
	assert.Regexp(t, regexp.MustCompile("node-2.*NotReady"), string(output))
}
