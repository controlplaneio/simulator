package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/runner"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Run(t *testing.T) {
	expected := readFixture("tf-help.txt")
	out, err := runner.Run("./", []string{}, "terraform", "help")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "out was nil")
	output := *out
	assert.Equal(t, output, expected)
}

func Test_Run_invalid_working_dir(t *testing.T) {
	wd := strings.Repeat("deadbeef", 1024)
	out, err := runner.Run(wd, []string{}, "terraform", "help")

	assert.NotNil(t, err, "Got no error")
	assert.Nil(t, out, "Got output")
	assert.Regexp(t, "^Error starting child process", err.Error())
}
