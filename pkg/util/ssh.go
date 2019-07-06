package util

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"time"
)

const (
	timeout = 10 * time.Minute
)

func agentSigners() ([]ssh.Signer, error) {
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}

	agent := agent.NewClient(sock)

	return agent.Signers()
}

func SSH(host string) error {
	port := "22"
	user := "root"

	signers, err := agentSigners()
	if err != nil {
		return err
	}
	// or get the signer from a private key file directly
	// signer, err := ssh.ParsePrivateKey(pemBytes)
	// if err != nil {
	//     log.Fatal(err)
	// }

	auths := []ssh.AuthMethod{ssh.PublicKeys(signers...)}

	// get host public key
	//hostKey := getHostKey(host)

	fmt.Printf("Connecting to %s\n", host)
	cfg := ssh.ClientConfig{
		User: user,
		Auth: auths,
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// verify host public key
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return StartInteractiveShell(&cfg, "tcp", host, port)
}

func StartInteractiveShell(sshConfig *ssh.ClientConfig, network string, host string, port string) error {
	var (
		session *ssh.Session
		conn    *ssh.Client
		err     error
	)
	if conn, err = getSSHConnection(sshConfig, network, host, port); err != nil {
		fmt.Printf("Failed to dial: %s", err)
		return err
	}

	if session, err = getSSHSession(conn); err != nil {
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

	if err = session.Shell(); err != nil {
		fmt.Printf("Failed to start interactive shell: %s", err)
		return err
	}
	return session.Wait()
}

func getSSHConnection(config *ssh.ClientConfig, network string, host string, port string) (*ssh.Client, error) {
	addr := host + ":" + port
	return ssh.Dial(network, addr, config)
}

func getSSHSession(clientConnection *ssh.Client) (*ssh.Session, error) {
	return clientConnection.NewSession()
}

// pty = pseudo terminal
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
