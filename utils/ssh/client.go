package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	bastion = "bastion"

	KeyAlgoED25519 = ssh.KeyAlgoED25519
)

func NewClient(sshConfigDir string, sshConfigName string, hostKeyAlg []string, timeout time.Duration) (*Client, error) {
	sshConfigFile, err := os.Open(filepath.Join(sshConfigDir, sshConfigName))
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}

	sshConfig, err := ssh_config.Decode(sshConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ssh config: %w", err)
	}

	user, err := sshConfig.Get(bastion, "User")
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}

	host, err := sshConfig.Get(bastion, "Hostname")
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}

	idFile, err := sshConfig.Get(bastion, "IdentityFile")
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}

	knownHostFile, err := sshConfig.Get(bastion, "UserKnownHostsFile")
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}

	client := &Client{
		user:           user,
		idFile:         filepath.Join(sshConfigDir, idFile),
		host:           host,
		hostKeyAlg:     hostKeyAlg,
		knownHostsFile: filepath.Join(sshConfigDir, knownHostFile),
		timeout:        timeout,
	}

	return client, nil
}

type Client struct {
	user           string
	idFile         string
	host           string
	hostKeyAlg     []string
	knownHostsFile string
	timeout        time.Duration
}

func (c Client) Execute(command string) ([]byte, error) {
	key, err := os.ReadFile(c.idFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ssh id file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh id file: %w", err)
	}

	hostKeyCallback, err := knownhosts.New(c.knownHostsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh known hosts file: %w", err)
	}

	sshClientConf := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyAlgorithms: c.hostKeyAlg,
		HostKeyCallback:   hostKeyCallback,
		Timeout:           c.timeout,
	}

	conn, err := ssh.Dial("tcp", net.JoinHostPort(c.host, "22"), sshClientConf)
	if err != nil {
		return nil, fmt.Errorf("failed to start ssh connection: %w", err)
	}
	defer func(conn *ssh.Client) {
		_ = conn.Close()
	}(conn)

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start ssh session: %w", err)
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to set ssh output: %w", err)
	}

	err = session.Start(command)
	if err != nil {
		return nil, fmt.Errorf("failed to run ssh command: %w", err)
	}

	output, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read ssh output: %w", err)
	}

	return output, nil
}
