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

// KeyScan runs ssh-keyscan silently against the provided bastion address. It returns a pointer to a string containing
// its buffered stdout or an error if any occurred
func KeyScan(bastion string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := RunSilently(wd, os.Environ(), "ssh-keyscan", "-H", bastion)
	return out, err
}

// GenerateKey runs ssh-keygen silently to create an SSH key with the same provided using preconfigured settings
// It returns a pointer to a string containing the buffered stdout or an error if any occurred
func GenerateKey(keyname string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	out, _, err := RunSilently(wd, os.Environ(), "ssh-keygen", "-f", keyname, "-t", "rsa", "-C", "''", "-N", "''")
	return out, err
}

// PublicKeyFile reads the public key at the path supplied and returns the ssh.AuthMethod to use or an error if any
// occurredj
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := Slurp(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey([]byte(*buffer))
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// GetAuthMethods tries to contact ssh-agent to get the AuthMethods and falls back to reading the keyfile directly
// in case of a missing SSH_AUTH_SOCK env var or an error dialing the unix socket
func GetAuthMethods() ([]ssh.AuthMethod, error) {
	// Check we have the ssh-agent AUTH SOCK and short circuit if we don't - just create a signer from the keyfile
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		fmt.Println("SSH_AUTH_SOCK was not set - falling back to id_rsa keyfile - this will fail if you have a passphrase")
		keypath, err := ExpandTilde("~/.ssh/id_rsa.pub")
		if err != nil {
			fmt.Println("Error reading id_rsa when falling back to key")

			return nil, err
		}

		return []ssh.AuthMethod{PublicKeyFile(*keypath)}, nil
	}

	// Try to get a signer from ssh-agent
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		// Fallback to keyfile if we failed to connect
		fmt.Println("Error dialing SSH_AUTH_SOCK - falling back to id_rsa key")
		keypath, err := ExpandTilde("~/.ssh/id_rsa.pub")
		if err != nil {
			fmt.Println("Error reading id_rsa when falling back to key")
			return nil, err
		}

		return []ssh.AuthMethod{PublicKeyFile(*keypath)}, nil
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
	// or get the signer from a private key file directly
	// signer, err := ssh.ParsePrivateKey(pemBytes)
	// if err != nil {
	//     log.Fatal(err)
	// }

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

	encodedkey, err := Base64PrivateKey("id_rsa")
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
