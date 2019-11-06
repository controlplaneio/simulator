package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

var noopLogger = zap.NewNop().Sugar()

var tfCommandArgumentsTests = []struct {
	command   string
	arguments []string
}{
	{"output", []string{"output", "-json"}},
	{"init", []string{"init", "-input=false", "--var-file=settings/bastion.tfvars"}},
	{"plan", []string{"plan", "-input=false", "--var-file=settings/bastion.tfvars"}},
	{"apply", []string{"apply", "-input=false", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
	{"destroy", []string{"destroy", "-input=false", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.command, func(t *testing.T) {
			assert.Equal(t, simulator.PrepareTfArgs(tt.command), tt.arguments)
		})
	}
}

func Test_Status(t *testing.T) {
	tfo, err := simulator.Status(noopLogger, fixture("noop-tf-dir"), "test", "latest")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfo, "Got no terraform output")
}

func Test_Create(t *testing.T) {
	err := simulator.Create(noopLogger, fixture("noop-tf-dir"), "test", "latest")

	assert.Nil(t, err)
}

func Test_Destroy(t *testing.T) {
	err := simulator.Destroy(noopLogger, fixture("noop-tf-dir"), "test", "test")

	assert.Nil(t, err)
}
