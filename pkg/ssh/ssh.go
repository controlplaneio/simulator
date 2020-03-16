package ssh

import (
	"github.com/kubernetes-simulator/simulator/pkg/progress"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/http"
	"os"
)

// GetAuthMethods tries to contact ssh-agent to get the AuthMethods and falls
// back to reading the keyfile directly in case of a missing SSH_AUTH_SOCK env
// var or an error dialing the unix socket
func GetAuthMethods(kp KeyPair) ([]ssh.AuthMethod, error) {
	keyFileAuth, err := kp.PrivateKey.ToAuthMethod()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting auth methods for private key")
	}

	return []ssh.AuthMethod{keyFileAuth}, nil
}

// SSH establishes an interactive Secure Shell session to the supplied host as
// user ubuntu and on port 22. SSH uses ssh-agent to get the key to use
func SSH(host string, kp KeyPair) error {
	port := "22"
	user := "ubuntu"

	auths, err := GetAuthMethods(kp)
	if err != nil {
		return errors.Wrap(err, "Error getting auth methods")
	}

	log.Printf("Connecting to %s\n", host)

	abspath, err := util.ExpandTilde(KnownHostsPath)
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

	return StartInteractiveSSHShell(&cfg, "tcp", host, port, kp)
}

// StartRemoteListener sets up a remote listener on the SSH connection
func StartRemoteListener(client *ssh.Client) {
	listener, err := client.Listen("tcp", "0.0.0.0:51234")
	if err != nil {
		log.Printf("Unable to start remote listener on SSH connection: %-v\n", err)
		return
	}

	handler := progress.NewHTTPHandler(progress.LocalStateProvider{})

	if err := http.Serve(listener, handler); err != nil {
		log.Printf("Unable to serve HTTP on the remote listener: %-v\n", err)

	}
}

// StartInteractiveSSHShell starts an interactive SSH shell with the supplied
// ClientConfig
func StartInteractiveSSHShell(sshConfig *ssh.ClientConfig, network string, host string, port string, kp KeyPair) error {
	var (
		session *ssh.Session
		conn    *ssh.Client
		err     error
	)

	f, err := os.OpenFile(util.MustExpandTilde("~/.kubesim/ssh-log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	addr := host + ":" + port
	if conn, err = ssh.Dial(network, addr, sshConfig); err != nil {
		log.Printf("Failed to dial: %s", err)
		return errors.Wrapf(err, "Error dialing %s", addr)
	}

	if session, err = conn.NewSession(); err != nil {
		log.Printf("Failed to create session: %s", err)
		return errors.Wrap(err, "Error establishing SSH session")
	}
	defer session.Close()

	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		// See this for more information http://www.linusakesson.net/programming/tty/
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return errors.Wrap(err, "Error setting stdin terminal to raw mode")
		}
		defer func() {
			if err := terminal.Restore(fileDescriptor, originalState); err != nil {
				// Something really bad happened
				panic(err)
			}
		}()
	}

	go StartRemoteListener(conn)

	if err = setupPty(fileDescriptor, session); err != nil {
		log.Printf("Failed to set up pseudo terminal: %s", err)
		return errors.Wrap(err, "Error setting up pseudo terminal")
	}

	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err = session.Setenv("BASE64_SSH_KEY", kp.PrivateKey.ToBase64()); err != nil {
		log.Printf("Failed to send SetEnv request: %s", err)
		return errors.Wrap(err, "Failed to send BASE_64_SSH_KEY env var using Setenv")
	}

	if err = session.Shell(); err != nil {
		log.Printf("Failed to start interactive shell: %s", err)
		return errors.Wrap(err, "Failed to start interactive shell")
	}

	return session.Wait()
}

// See http://www.tldp.org/HOWTO/Text-Terminal-HOWTO-7.html#ss7.2 for more info
// on pseudo terminals
func setupPty(stdinFd int, session *ssh.Session) error {
	// https://tools.ietf.org/html/rfc4254#section-8 for more information about
	// terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing of characters as you type
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	termWidth, termHeight, err := terminal.GetSize(stdinFd)
	if err != nil {
		return errors.Wrap(err, "Error getting size of stdin terminal")
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		session.Close()
		return errors.Wrap(err, "Error sending pty request for an xterm over ssh session")
	}

	return nil
}
