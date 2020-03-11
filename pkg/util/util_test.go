package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixture(name string) string {
	return "../../test/fixtures/" + name
}

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

func Test_IsStringInSlice(t *testing.T) {
	b, err := util.IsStringInSlice("test", []string{"not-test", "test"})

	assert.Nil(t, err, "Got an error")
	assert.True(t, b, "String was not in slice")

	b, err = util.IsStringInSlice("test", []string{"not-test", "also-not-test"})
	assert.NotNil(t, err, "Got an error")
	assert.False(t, b, "String was present in slice")
}

func Test_ExpandTilde(t *testing.T) {
	p, err := util.ExpandTilde("~/.")
	require.NoError(t, err)

	assert.Regexp(t, `^/(home|Users)/([^/]+)$`, *p)

	// Call ExpandTilde again to exercise the cache
	p2, err := util.ExpandTilde("~/.")
	require.Nil(t, err, "Got an error resolving tilde")
	assert.Equal(t, *p, *p2, "Cached version differed")

	p3, err := util.ExpandTilde("fail")
	require.NotNil(t, err, "Didn't get an error when path didn't start with tilde slash")
	assert.Nil(t, p3, "Got a path when path didn't start with tilde slash")
	assert.Regexp(t, `^Path was empty or did not start with a tilde and a slash:`, err.Error())

	p4, err := util.ExpandTilde("")
	require.NotNil(t, err, "Didn't get an error for empty path")
	assert.Nil(t, p4, "Got a path when resolving an empty path")
	assert.Regexp(t, `^Path was empty or did not start with a tilde and a slash:`, err.Error())
}

func Test_Slurp(t *testing.T) {
	contents, err := util.Slurp(fixture("tf-help.txt"))
	assert.Nil(t, err, "Got an error")
	assert.NotEmpty(t, contents, "Got empty file contents")
}

func Test_MustSlurp(t *testing.T) {
	contents := util.MustSlurp(fixture("tf-help.txt"))
	assert.NotEmpty(t, contents, "Got empty file contents")
}

func Test_EnsureFile(t *testing.T) {
	written, err := util.EnsureFile(fixture("tf-help.txt"), "testing")
	assert.Nil(t, err, "Got an error")
	assert.False(t, written, "Wrote the file")

	util.MustRemove(fixture(".ignored"))
	written, err = util.EnsureFile(fixture(".ignored"), "testing")
	assert.Nil(t, err, "Got an error")
	assert.True(t, written, "Didn't write the file")
}

func Test_MustRemove(t *testing.T) {
	util.MustRemove("./non-existent-file")
}

func Test_OverwriteFile(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		deletefile string
		content    string
	}{
		{
			name:       "samedir",
			filename:   "test.yaml",
			deletefile: "test.yaml",
			content:    "hello",
		},
		{
			name:       "nestdir",
			filename:   filepath.Join("testdir", "1", "test.yaml"),
			deletefile: "./testdir",
			content:    "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.OverwriteFile(tt.filename, tt.content)
			assert.NoError(t, err, tt.name)
			defer func() {
				err = os.RemoveAll(tt.deletefile)
				assert.NoError(t, err, tt.name)
			}()
			bytes, err := ioutil.ReadFile(tt.filename)
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.content, string(bytes), tt.name)
		})
	}
}
