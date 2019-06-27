package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func checkWorkingDir(wd string) (string, error) {
	debug("Checking working dir")
	absPath, err := filepath.Abs(wd)
	if err != nil {
		return "", errors.Wrapf(err, "Error resolving working dir %s", wd)
	}

	return absPath, nil
}

// Run runs a child process and returns its buffer stdout.  Run also tees the output to stdout of this process, `env` will
// be appended to the current environment.  `wd` is the working directory for the child
func Run(wd string, env []string, cmd string, args ...string) (*string, error) {

	debug("Preparing to run: ", cmd, args)
	child := exec.Command(cmd, args...)

	child.Env = append(os.Environ(), env...)

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir, err := checkWorkingDir(wd)
	if err != nil {
		return nil, err
	}

	debug("Setting child working directory to ", dir)
	child.Dir = dir

	// Copy child stdout to stdout but also into a buffer to be returned
	var buf bytes.Buffer
	tee := io.TeeReader(childOut, &buf)

	debug("Running child")
	err = child.Start()
	if err != nil {
		debug("Error starting child process: ", err)
		return nil, errors.Wrapf(err, "Error starting child process")
	}

	io.Copy(os.Stdout, tee)
	io.Copy(os.Stderr, childErr)

	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		debug("Error waiting for child process", err)
		return nil, err
	}

	out := string(buf.Bytes())
	return &out, nil
}
