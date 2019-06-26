package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tfCommandArgumentsTests = []struct {
	command   string
	arguments []string
}{
	{"init", []string{"init", "--var-file=settings/bastion.tfvars"}},
	{"output", []string{"output", "-json"}},
	{"plan", []string{"plan", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
	{"apply", []string{"apply", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
	{"destroy", []string{"destroy", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareArguments(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.command, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, runner.PrepareArguments(tt.command), tt.arguments)
		})
	}
}
