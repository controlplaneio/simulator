package simulator_test

import (
	sim "github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

var noopLogger = zap.NewNop().Sugar()

var tfCommandArgumentsTests = []struct {
	prepArgs  []string
	arguments []string
}{
	{[]string{"output", "test-bucket"}, []string{"output", "-json"}},
	{[]string{"init", "test-bucket"}, []string{"init", "-input=false", "--var-file=~/.kubesim/settings/bastion.tfvars", "-backend-config=bucket=test-bucket"}},
	{[]string{"plan", "test-bucket"}, []string{"plan", "-input=false", "--var-file=~/.kubesim/settings/bastion.tfvars"}},
	{[]string{"apply", "test-bucket"}, []string{"apply", "-input=false", "--var-file=~/.kubesim/settings/bastion.tfvars", "-auto-approve"}},
	{[]string{"destroy", "test-bucket"}, []string{"destroy", "-input=false", "--var-file=~/.kubesim/settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.prepArgs[0], func(t *testing.T) {
			assert.Equal(t, sim.PrepareTfArgs(tt.prepArgs[0], tt.prepArgs[1]), tt.arguments)
		})
	}
}
	
func Test_Status(t *testing.T) {
	
	simulator := sim.NewSimulator(
		sim.WithLogger(noopLogger),
		sim.WithTfDir(fixture("noop-tf-dir")),
		sim.WithScenariosDir("test"),
		sim.WithAttackTag("latest"),
		sim.WithBucketName("test"),
		sim.WithTfVarsDir(fixture("noop-tf-dir")))

	tfo, err := simulator.Status()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfo, "Got no terraform output")
}

// func Test_Create(t *testing.T) {
	
// 	simulator := sim.NewSimulator(
// 		sim.WithLogger(noopLogger),
// 		sim.WithTfDir(fixture("noop-tf-dir")),
// 		sim.WithScenariosDir("test"),
// 		sim.WithAttackTag("latest"),
// 		sim.WithBucketName("test"),
// 		sim.WithTfVarsDir(fixture("noop-tf-dir")))

// 	err := simulator.Create()

// 	assert.Nil(t, err)
// }

func Test_Destroy(t *testing.T) {
	err := sim.Destroy(noopLogger, fixture("noop-tf-dir"), "test", "test", "test"))

	assert.Nil(t, err)
}
