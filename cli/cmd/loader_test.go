package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_loadScenarios(t *testing.T) {
	scenarios, err := loadScenarios("../../simulation-scripts/")

	assert.Nil(t, err)
	assert.NotEqual(t, len(scenarios), 0, "Returned no scenarios")
}

func Test_manifestPath_default(t *testing.T) {
	p := manifestPath()
	assert.Equal(t, p, "../simulation-scripts/")
}

func Test_manifestPath_custom(t *testing.T) {
	os.Setenv("SIMULATOR_MANIFEST_PATH", "/some/path")
	p := manifestPath()
	assert.Equal(t, p, "/some/path", "manifestPath did not set custom path")
}

// TODO: Test bad scenario manifests
