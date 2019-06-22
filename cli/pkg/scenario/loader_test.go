package scenario_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/scenario"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_loadScenarios(t *testing.T) {
	manifest, err := scenario.LoadManifest("../../../simulation-scripts/")

	assert.Nil(t, err)
	assert.NotEqual(t, len(manifest.Scenarios), 0, "Returned no scenarios")
}

func Test_manifestPath_default(t *testing.T) {
	p := scenario.ManifestPath()
	assert.Equal(t, p, "../simulation-scripts/")
}

func Test_manifestPath_custom(t *testing.T) {
	os.Setenv("SIMULATOR_MANIFEST_PATH", "/some/path")
	p := scenario.ManifestPath()
	assert.Equal(t, p, "/some/path", "manifestPath did not set custom path")
}

// TODO: Test bad scenario manifests
