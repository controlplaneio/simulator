package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, simulator.PrepareTfArgs(tt.command), tt.arguments)
		})
	}
}

func Test_Status(t *testing.T) {
	tfo, err := simulator.Status(fixture("noop-tf-dir"), "test")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfo, "Got no terraform output")
}

func Test_Create(t *testing.T) {
	err := simulator.Create(fixture("noop-tf-dir"), "test")

	assert.Nil(t, err)
}

func Test_Destroy(t *testing.T) {
	err := simulator.Destroy(fixture("noop-tf-dir"), "test")

	assert.Nil(t, err)
}
