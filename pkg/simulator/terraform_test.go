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
	{[]string{"output", "test-bucket", fixture("noop-tf-dir")}, []string{"output", "-json"}},
	{[]string{"init", "test-bucket", fixture("noop-tf-dir")}, []string{"init", "-input=false", "--var-file=/go/src/github.com/controlplaneio/simulator-standalone/pkg/simulator/settings/bastion.tfvars", "-backend-config=bucket=test-bucket"}},
	{[]string{"plan", "test-bucket", fixture("noop-tf-dir")}, []string{"plan", "-input=false", "--var-file=/go/src/github.com/controlplaneio/simulator-standalone/pkg/simulator/settings/bastion.tfvars"}},
	{[]string{"apply", "test-bucket", fixture("noop-tf-dir")}, []string{"apply", "-input=false", "--var-file=/go/src/github.com/controlplaneio/simulator-standalone/pkg/simulator/settings/bastion.tfvars", "-auto-approve"}},
	{[]string{"destroy", "test-bucket", fixture("noop-tf-dir")}, []string{"destroy", "-input=false", "--var-file=/go/src/github.com/controlplaneio/simulator-standalone/pkg/simulator/settings/bastion.tfvars", "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.prepArgs[0], func(t *testing.T) {
			//TODO: put in tfVarsDir arg, maybe above
			assert.Equal(t, sim.PrepareTfArgs(tt.prepArgs[0], tt.prepArgs[1], tt.prepArgs[2]), tt.arguments)
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

func Test_Create(t *testing.T) {
	
	simulator := sim.NewSimulator(
		sim.WithLogger(noopLogger),
		sim.WithTfDir(fixture("noop-tf-dir")),
		sim.WithScenariosDir("test"),
		sim.WithAttackTag("latest"),
		sim.WithBucketName("test"),
		sim.WithTfVarsDir(fixture("noop-tf-dir")))

	err := simulator.Create()

    // pwd, err := os.Getwd()
    // if err != nil {
    //     fmt.Println(err)
    //     os.Exit(1)
    // }
	// fmt.Println(pwd)
	
	// assert.Nil(t, pwd)

	// file, err := os.Open("/go/src/github.com/controlplaneio/simulator-standalone/pkg/simulator/test/fixtures/noop-tf-dir/settings/bastion.tfvars")
    // if err != nil {
    //     log.Fatal(err)
    // }

	// n := fmt.Sprint(file)
	// assert.Nil(t, n)
	// assert.Nil(t, err)
	
	assert.Nil(t, err)
}

func Test_Destroy(t *testing.T) {
	// err := sim.Destroy(noopLogger, fixture("noop-tf-dir"), "test", "test", fixture("noop-tf-dir"))

	// assert.Nil(t, err)
}
