package runner

import (
	"fmt"
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
