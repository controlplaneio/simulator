package runner

import (
	"fmt"
	"io/ioutil"
	"os"
)

func debug(msg ...interface{}) {
	fmt.Println(msg...)
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
	file, err := os.Open(path)
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
