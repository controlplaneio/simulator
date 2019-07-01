package util_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_FileExists_current_file(t *testing.T) {
	exists, err := util.FileExists("./util_test.go")

	assert.Nil(t, err, "Got an error")
	assert.True(t, exists, "Didn't return true for current file")
}

func Test_FileExists_garbage(t *testing.T) {
	exists, err := util.FileExists("./non-existent.garbage")

	assert.Nil(t, err, "Got an error")
	assert.False(t, exists, "Didn't return false for garbage")
}

func Test_Slurp_current_file(t *testing.T) {
	contents, err := util.Slurp("./util_test.go")

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, contents, "Didn't return file contents")
}

func Test_DetectPublicIP(t *testing.T) {
	ip, err := util.DetectPublicIP()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, ip, "Got no IP address")
}

func Test_EnvOrDefault(t *testing.T) {
	key := "SIMULATOR_TEST_" + string(time.Now().Unix())
	defaulted := util.EnvOrDefault(key, "setting")
	assert.Equal(t, defaulted, "setting", "Did not return default")

	os.Setenv(key, "custom")
	val := util.EnvOrDefault(key, "custom")
	assert.Equal(t, val, "custom", "Did not read env var")
}

func Test_ExpandTilde(t *testing.T) {
	p, err := util.ExpandTilde("~/.")
	assert.Nil(t, err, "Got an error")

	assert.Regexp(t, `^/home/([^/]+)$`, *p)

	// Call ExpandTilde again to exercise the cache
	p2, err := util.ExpandTilde("~/.")
	assert.Equal(t, *p, *p2, "Cached version differed")
}
