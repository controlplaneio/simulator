package ssh

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
	"time"
)

const (
	timeout = 10 * time.Minute
)

// GetAuthMethods tries to contact ssh-agent to get the AuthMethods and falls back to reading the keyfile directly
// in case of a missing SSH_AUTH_SOCK env var or an error dialing the unix socket
func GetAuthMethods() ([]ssh.AuthMethod, error) {
	// Check we have the ssh-agent AUTH SOCK and short circuit if we don't - just create a signer from the keyfile
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		keyFileAuth, err := PrivateKeyFile()
		if err != nil {
			return nil, err
		}

		fmt.Println("KeyFile")
		fmt.Printf("%+v", keyFileAuth)
		return []ssh.AuthMethod{keyFileAuth}, nil
	}

	// Try to get a signer from ssh-agent
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		// Fallback to keyfile if we failed to connect
		keyFileAuth, err := PrivateKeyFile()
		if err != nil {
			return nil, err
		}

		return []ssh.AuthMethod{keyFileAuth}, nil
	}

	// Use the signers from ssh-agent
	agent := agent.NewClient(sock)

	signers, err := agent.Signers()
	if err != nil {
		return nil, err
	}

	return []ssh.AuthMethod{ssh.PublicKeys(signers...)}, nil
}

// SSH establishes an interactive Secure Shell session to the supplied host as user ubuntu and on port 22. SSH uses
// ssh-agent to get the key to use
func SSH(host string) error {
	port := "22"
	user := "ubuntu"

	auths, err := GetAuthMethods()
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to %s\n", host)

	abspath, err := util.ExpandTilde(SSHKnownHostsPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving known_hosts path")
	}

	knownHostsCallback, err := knownhosts.New(*abspath)
	if err != nil {
		return errors.Wrap(err, "Error configuring ssh client to use known_hosts file")
	}

	cfg := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: knownHostsCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
	}

	return StartInteractiveSSHShell(&cfg, "tcp", host, port)
}

// StartInteractiveSSHShell starts an interactive SSH shell with the supplied ClientConfig
func StartInteractiveSSHShell(sshConfig *ssh.ClientConfig, network string, host string, port string) error {
	var (
		session *ssh.Session
		conn    *ssh.Client
		err     error
	)

	addr := host + ":" + port
	if conn, err = ssh.Dial(network, addr, sshConfig); err != nil {
		fmt.Printf("Failed to dial: %s", err)
		return err
	}

	encodedkey, err := Base64PrivateKey(PrivateKeyPath)
	if err != nil {
		return err
	}

	if session, err = conn.NewSession(); err != nil {
		fmt.Printf("Failed to create session: %s", err)
		return err
	}
	defer session.Close()

	if err = setupPty(session); err != nil {
		fmt.Printf("Failed to set up pseudo terminal: %s", err)
		return err
	}

	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err = session.Setenv("BASE64_SSH_KEY", *encodedkey); err != nil {
		fmt.Printf("Failed to send SetEnv request: %s", err)
		return err
	}

	if err = session.Shell(); err != nil {
		fmt.Printf("Failed to start interactive shell: %s", err)
		return err
	}

	return session.Wait()
}

func setupPty(session *ssh.Session) error {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		fmt.Printf("request for pseudo terminal failed: %s", err)
		return err
	}
	return nil
}
