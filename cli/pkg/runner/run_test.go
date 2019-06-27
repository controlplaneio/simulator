package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/runner"
	"github.com/stretchr/testify/assert"
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
