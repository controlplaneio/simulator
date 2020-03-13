package simulator_test

import (
	sim "github.com/kubernetes-simulator/simulator/pkg/simulator"
	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
	"os"
)

var pwd, _ = os.Getwd()
var testVarFileArg = "--var-file=" + pwd + "/" + fixture("noop-tf-dir") + "/settings/bastion.tfvars"
var logger = logrus.New()

var tfCommandArgumentsTests = []struct {
	prepArgs  []string
	arguments []string
}{
	{[]string{"output"}, []string{"output", "-json"}},
	{[]string{"init"}, []string{"init", "-input=false", testVarFileArg, "-backend-config=bucket=test-bucket"}},
	{[]string{"plan"}, []string{"plan", "-input=false", testVarFileArg}},
	{[]string{"apply"}, []string{"apply", "-input=false", testVarFileArg, "-auto-approve"}},
	{[]string{"destroy"}, []string{"destroy", "-input=false", testVarFileArg, "-auto-approve"}},
}

func Test_PrepareTfArgs(t *testing.T) {
	pwd, _ := os.Getwd()
	logger.Out = ioutil.Discard
	simulator := sim.NewSimulator(
		sim.WithLogger(logger),
		sim.WithBucketName("test-bucket"),
		sim.WithoutIPDetection(false),
		sim.WithTfVarsDir(pwd+"/"+fixture("noop-tf-dir")))

	for _, tt := range tfCommandArgumentsTests {
		t.Run("Test arguments for "+tt.prepArgs[0], func(t *testing.T) {
			assert.Equal(t, simulator.PrepareTfArgs(tt.prepArgs[0]), tt.arguments)
		})
	}
}

func Test_Status(t *testing.T) {
	os.Remove(util.MustExpandTilde(ssh.PublicKeyPath))
	os.Remove(util.MustExpandTilde(ssh.PrivateKeyPath))
	pwd, _ := os.Getwd()
	logger.Out = ioutil.Discard
	simulator := sim.NewSimulator(
		sim.WithLogger(logger),
		sim.WithTfDir(fixture("noop-tf-dir")),
		sim.WithScenariosDir("test"),
		sim.WithAttackTag("latest"),
		sim.WithBucketName("test"),
		sim.WithoutIPDetection(false),
		sim.WithTfVarsDir(pwd+"/"+fixture("noop-tf-dir")))

	tfo, err := simulator.Status()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfo, "Got no terraform output")
}

func Test_Create(t *testing.T) {

	pwd, _ := os.Getwd()
	logger.Out = ioutil.Discard
	simulator := sim.NewSimulator(
		sim.WithLogger(logger),
		sim.WithTfDir(fixture("noop-tf-dir")),
		sim.WithScenariosDir("test"),
		sim.WithAttackTag("latest"),
		sim.WithBucketName("test"),
		sim.WithTfVarsDir(pwd+"/"+fixture("noop-tf-dir")))

	err := simulator.Create()
	assert.Nil(t, err)
}

func Test_Destroy(t *testing.T) {

	pwd, _ := os.Getwd()
	logger.Out = ioutil.Discard
	simulator := sim.NewSimulator(
		sim.WithLogger(logger),
		sim.WithTfDir(fixture("noop-tf-dir")),
		sim.WithAttackTag("latest"),
		sim.WithBucketName("test"),
		sim.WithTfVarsDir(pwd+"/"+fixture("noop-tf-dir")))

	err := simulator.Destroy()

	assert.Nil(t, err)
}
