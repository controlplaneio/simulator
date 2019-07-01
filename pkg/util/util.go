package util

import (
	"fmt"
	"github.com/glendc/go-external-ip"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

// Debug prints a debug message to stdout
func Debug(msg ...interface{}) {
	fmt.Println(msg...)
}

// DetectPublicIP detects your public IP address
func DetectPublicIP() (*string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, err
	}

	output := ip.String()
	return &output, nil
}

var homedirCache string
var cacheLock sync.RWMutex

// ExpandTilde returns the fully qualified path to a file in the user's home directory. I.E. it expands a path beginning with
// `~/`) and checks the file exists. ExpandTilde will cache the user's home directory to amortise the cost of the syscall
func ExpandTilde(path string) (*string, error) {
	if len(path) == 0 || path[0:2] != "~/" {
		return nil, errors.Errorf(`Path was empty or did not start with a tilde and a slash: "%s"`, path)
	}

	// discard ~/
	path = path[2:]

	var homedir string

	// Lock and read the cache to see if we already resolved the current user's home directory
	cacheLock.RLock()
	homedir = homedirCache
	cacheLock.RUnlock()
	if homedir == "" {
		// Take a write lock to update the cache
		cacheLock.Lock()
		defer cacheLock.Unlock()

		usr, err := user.Current()
		if err != nil {
			return nil, errors.Wrapf(err, "Error finding %s in home", path)
		}

		homedir = usr.HomeDir
	}

	p := filepath.Join(homedir, path)
	exists, err := FileExists(p)
	if err != nil {
		return nil, errors.Wrapf(err, "Error checking %s exists", p)
	}
	if !exists {
		return nil, errors.Errorf("No file found at %s", p)
	}

	return &p, nil
}

// FileExists checks whether a path exists
func FileExists(path string) (bool, error) {
	Debug("Stating", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Slurp reads an entire file into a string in one operation and returns a pointer to the file's content or an error if
// any.  Similar to `ioutil.ReadFile` but it calls `filepath.Abs` first which cleans the path and resolves relative
// paths from the working directory.
//
// Note that this is slightly less efficient for zero-length files than `ioutil.Readfile` as it uses the default
// read buffer size of `bytes.MinRead` internally
func Slurp(path string) (*string, error) {
	fp, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	output := string(b)
	return &output, nil
}

// MustSlurp is the panicky counterpart to Slurp. MustSlurp reads an entire file into a string in one operation and
// returns the contents or panics if it encouters and error
func MustSlurp(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// MustRemove removes a file or empty directory.  MustRemove will ignore an error if the path doesn't exist or panic for
// any other error
func MustRemove(path string) {
	err := os.Remove(path)
	if err != nil && os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}

	return
}

// EnsureFile checks a file exists and writes the supplied contents if not.  returns a boolean indicating whether it
// wrote a file or not and any error
func EnsureFile(path, contents string) (bool, error) {
	exists, err := FileExists(path)
	if err != nil || exists {
		return false, err
	}

	file, err := os.Create(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = io.WriteString(file, contents)
	if err != nil {
		return false, err
	}

	err = file.Sync()
	if err != nil {
		return true, err
	}

	return true, nil
}

// EnvOrDefault tries to read an environment variable with the supplied key and returns its value.  EnvOrDefault returns
// a default value if it is empty or unset
func EnvOrDefault(key, def string) string {
	var d = os.Getenv(key)
	if d == "" {
		d = def
	}

	return d
}
