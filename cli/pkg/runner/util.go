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

// DefaultPublicKeyPath returns the path to `id_rsa.pub`
func DefaultPublicKeyPath() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding default public key")
	}

	home := usr.HomeDir
	keypath := filepath.Join(home, ".ssh/id_rsa.pub")
	exists, err := FileExists(keypath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error checking default public key exists at %s", keypath)
	}
	if !exists {
		return nil, errors.Errorf("No default public key found at %s", keypath)
	}

	return &keypath, nil
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

// EnsureFile checks a file exists and writes the supplied contents if not
func EnsureFile(path, contents string) error {
	exists, err := FileExists(path)
	if err != nil || exists {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, contents)
	if err != nil {
		return err
	}

	return file.Sync()
}

// EnvOrDefault tries to read the key and returns a default value if it is empty
func EnvOrDefault(key, def string) string {
	var d = os.Getenv(key)
	if d == "" {
		d = def
	}

	return d
}
