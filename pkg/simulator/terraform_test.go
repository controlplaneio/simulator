package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

var noopLogger = zap.NewNop().Sugar()

var tfCommandArgumentsTests = []struct {
	prepArgs   []string
	arguments []string
}{
	{[]string{"output", "test-bucket"}, []string{"output", "-json"}},
	{[]string{"init", "test-bucket"}, []string{"init", "-input=false", "--var-file=settings/bastion.tfvars", "-backend-config=bucket=test-bucket"}},
	{[]string{"plan", "test-bucket"}, []string{"plan", "-input=false", "--var-file=settings/bastion.tfvars"}},
	{[]string{"apply", "test-bucket"}, []string{"apply", "-input=false", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
	{[]string{"destroy", "test-bucket"}, []string{"destroy", "-input=false", "--var-file=settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.prepArgs[0], func(t *testing.T) {
			assert.Equal(t, simulator.PrepareTfArgs(tt.prepArgs[0], tt.prepArgs[1]), tt.arguments)
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
