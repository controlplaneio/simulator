package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FileExists_current_file(t *testing.T) {
	exists, err := runner.FileExists("./util_test.go")

	assert.Nil(t, err, "Got an error")
	assert.True(t, exists, "Didn't return true for current file")
}

func Test_FileExists_garbage(t *testing.T) {
	exists, err := runner.FileExists("./non-existent.garbage")

	assert.Nil(t, err, "Got an error")
	assert.False(t, exists, "Didn't return false for garbage")
}

func Test_ReadFile_current_file(t *testing.T) {
	contents, err := runner.ReadFile("./util_test.go")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, contents, "Didn't return file contents")
}

func Test_DetectPublicIP(t *testing.T) {
	ip, err := runner.DetectPublicIP()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, ip, "Got no IP address")
}
