package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/runner"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var tfCommandArgumentsTests = []struct {
	command   string
	arguments []string
}{
	{"init", []string{"init", "--var-file=settings/bastion.tfvars"}},
	{"output", []string{"output", "-json"}},
	{"plan", []string{"plan", "--var-file=settings/bastion.tfvars"}},
	{"apply", []string{"apply", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
	{"destroy", []string{"destroy", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.command, func(t *testing.T) {
			assert.Equal(t, runner.PrepareTfArgs(tt.command), tt.arguments)
		})
	}
}

func Test_Terraform(t *testing.T) {
	os.Setenv("SIMULATOR_TF_DIR", "../../../terraform/deployments/AwsSimulatorStandalone")
	expected := readFixture("tf-help.txt")
	out, err := runner.Terraform("help")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "out was nil")
	output := *out
	assert.Equal(t, output, expected)
}
