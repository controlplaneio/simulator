package runner

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

func debug(msg ...interface{}) {
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

// Home returns the fully qualified path to a file in the user's home directory. I.E. it expands a path beginning with
// `~`) and checks the file exists
func Home(path string) (*string, error) {
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
	debug("Stating", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ReadFile return a pointer to a string with the file's content
func ReadFile(path string) (*string, error) {
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

// EnvOrDefault tries to read the key and returns a default value if it is empty
func EnvOrDefault(key, def string) string {
	var d = os.Getenv(key)
	if d == "" {
		d = def
	}

	return d
}
