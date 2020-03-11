package childminder_test

import (
	"strings"
	"testing"

	"github.com/kubernetes-simulator/simulator/pkg/childminder"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func fixture(name string) string {
	return "../../test/fixtures/" + name
}

func Test_Run(t *testing.T) {
	expected := util.MustSlurp(fixture("tf-help.txt"))
	logger := logrus.New()
	cm := childminder.NewChildMinder(logger, "./", []string{}, "terraform", "help")
	out, err := cm.Run()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "out was nil")
	output := *out
	assert.Equal(t, output, expected)
}

func Test_Run_invalid_working_dir(t *testing.T) {
	wd := strings.Repeat("deadbeef", 1024)
	logger := logrus.New()
	cm := childminder.NewChildMinder(logger, wd, []string{}, "terraform", "help")
	out, err := cm.Run()

	assert.NotNil(t, err, "Got no error")
	assert.Nil(t, out, "Got output")
	assert.Regexp(t, "^Error starting child process", err.Error())
}

func Test_Run_silently(t *testing.T) {
	expected := util.MustSlurp(fixture("tf-help.txt"))
	logger := logrus.New()
	cm := childminder.NewChildMinder(logger, "./", []string{}, "terraform", "help")
	out, _, err := cm.RunSilently()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, out, "out was nil")
	output := *out
	assert.Equal(t, output, expected)
}
