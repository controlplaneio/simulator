package childminder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ChildMinder represents a child process to be managed
type ChildMinder struct {
	// Logger is the logger the simulator will use
	Logger *logrus.Logger
	// The working directory of the chlid process
	WorkingDir string
	// Environment is the variables to append to the parent's environment when
	// running the child process
	Environment []string
	// CommandPath is the path to the command to run
	CommandPath string
	// CommmandArguments are the arguments to supply to the command
	CommandArguments []string
}

// NewChildMinder creates an instance of a ChildMinder
func NewChildMinder(
	logger *logrus.Logger,
	wd string,
	env []string,
	cmd string,
	args ...string) *ChildMinder {

	cm := ChildMinder{}
	cm.Logger = logger
	cm.WorkingDir = wd
	cm.Environment = env
	cm.CommandPath = cmd
	cm.CommandArguments = args

	return &cm
}

// MustResolve returns an absolute path or panics if the underlying sys call
// fails
func MustResolve(wd string) string {
	absPath, err := filepath.Abs(wd)
	if err != nil {
		panic(err)
	}

	return absPath
}

// ForwardStdOut takes a child process's stdout pipe and reads it
// line by line passing the output through the ChildMinder's logger
func (cm *ChildMinder) ForwardStdOut(stdoutPipe io.Reader, wg *sync.WaitGroup) {
	stdoutReader := bufio.NewReader(stdoutPipe)
	for {
		line, err := stdoutReader.ReadString('\n')
		if len(line) > 0 {
			cm.Logger.Infof(fmt.Sprintf("[%s] %s", cm.CommandPath, line))
		}

		if err == io.EOF {
			cm.Logger.WithFields(logrus.Fields{
				"Command": cm.CommandPath,
				"Args":    cm.CommandArguments,
				"Error":   err,
			}).Debug("Stdout pipe has been closed")
			wg.Done()
			return
		}

		if err != nil {
			cm.Logger.WithFields(logrus.Fields{
				"Command": cm.CommandPath,
				"Args":    cm.CommandArguments,
				"Error":   err,
			}).Error("Error reading stdout pipe")
			wg.Done()
			return
		}
	}
}

// ForwardStdErr takes a child process's stderr pipe and reads it
// line by line passing the output through the ChildMinder's logger
func (cm *ChildMinder) ForwardStdErr(stderrPipe io.Reader, wg *sync.WaitGroup) {
	stderrReader := bufio.NewReader(stderrPipe)
	for {
		line, err := stderrReader.ReadString('\n')
		if len(line) > 0 {
			cm.Logger.WithFields(logrus.Fields{
				"Command": cm.CommandPath,
				"Args":    cm.CommandArguments,
			}).Error(line)
		}

		if err == io.EOF {
			cm.Logger.WithFields(logrus.Fields{
				"Command": cm.CommandPath,
				"Args":    cm.CommandArguments,
				"Error":   err,
			}).Debug("Stderr pipe has been closed")
			wg.Done()
			return
		}

		if err != nil {
			cm.Logger.WithFields(logrus.Fields{
				"Command": cm.CommandPath,
				"Args":    cm.CommandArguments,
				"Error":   err,
			}).Error("Error reading stderr pipe")
			wg.Done()
			return
		}
	}
}

// Run runs a child process and returns its buffer stdout.  Run also tees the
// output to stdout of this process, `env` will be appended to the current
// environment.  `wd` is the working directory for the child
func (cm *ChildMinder) Run() (*string, error) {
	child := exec.Command(cm.CommandPath, cm.CommandArguments...)

	child.Env = append(os.Environ(), cm.Environment...)

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir := MustResolve(cm.WorkingDir)

	child.Dir = dir

	// Copy child stdout to stdout but also into a buffer to be returned
	var buf bytes.Buffer
	tee := io.TeeReader(childOut, &buf)

	var wg sync.WaitGroup

	go cm.ForwardStdOut(tee, &wg)
	go cm.ForwardStdErr(childErr, &wg)
	wg.Add(2)

	err := child.Start()
	if err != nil {
		return nil, errors.Wrapf(err, "Error starting child process %s", cm.CommandPath)
	}

	wg.Wait()
	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		return nil, errors.Wrapf(err, "Error waiting for child process %s", cm.CommandPath)
	}

	out := buf.String()
	return &out, nil
}

// RunSilently runs a sub command silently
func (cm *ChildMinder) RunSilently() (*string, *string, error) {
	child := exec.Command(cm.CommandPath, cm.CommandArguments...)

	child.Env = append(os.Environ(), cm.Environment...)

	var outBuf, errBuf bytes.Buffer
	child.Stdout = bufio.NewWriter(&outBuf)
	child.Stderr = bufio.NewWriter(&errBuf)
	dir := MustResolve(cm.WorkingDir)

	child.Dir = dir

	err := child.Start()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error starting child process %s", cm.CommandPath)
	}

	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		childOut := outBuf.String()
		childErr := errBuf.String()
		return &childOut, &childErr, err
	}

	childErr := errBuf.String()
	childOut := outBuf.String()

	return &childOut, &childErr, nil
}
