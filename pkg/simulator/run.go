package simulator

import (
	"bytes"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func wdMust(wd string) string {
	// only errors if syscall to get parent's wd fails in which case something
	// really bad happened
	absPath, err := filepath.Abs(wd)
	if err != nil {
		panic(err)
	}

	return absPath
}

// Run runs a child process and returns its buffer stdout.  Run also tees the output to stdout of this process, `env` will
// be appended to the current environment.  `wd` is the working directory for the child
func Run(wd string, env []string, cmd string, args ...string) (*string, error) {

	util.Debug("Preparing to run: ", cmd, args)
	child := exec.Command(cmd, args...)

	child.Env = append(os.Environ(), env...)

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir := wdMust(wd)

	util.Debug("Setting child working directory to ", dir)
	child.Dir = dir

	// Copy child stdout to stdout but also into a buffer to be returned
	var buf bytes.Buffer
	tee := io.TeeReader(childOut, &buf)

	util.Debug("Running child")
	err := child.Start()
	if err != nil {
		util.Debug("Error starting child process: ", err)
		return nil, errors.Wrapf(err, "Error starting child process")
	}

	io.Copy(os.Stdout, tee)
	io.Copy(os.Stderr, childErr)

	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		util.Debug("Error waiting for child process", err)
		return nil, err
	}

	out := string(buf.Bytes())
	return &out, nil
}
