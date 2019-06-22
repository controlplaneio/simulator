package scenario_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/scenario"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LoadManifest(t *testing.T) {
	manifest, err := scenario.LoadManifest("../../../simulation-scripts/")

	assert.Nil(t, err)
	assert.NotEqual(t, len(manifest.Scenarios), 0, "Returned no scenarios")
}

func Test_ManifestPath_default(t *testing.T) {
	p := scenario.ManifestPath()
	assert.Equal(t, p, "../simulation-scripts/")
}

func Test_ManifestPath_custom(t *testing.T) {
	os.Setenv("SIMULATOR_MANIFEST_PATH", "/some/path")
	p := scenario.ManifestPath()
	assert.Equal(t, p, "/some/path", "manifestPath did not set custom path")
}

func fixture(name string) string {
	return "../../test/fixtures/" + name
}

func Test_LoadManifest_missing_manifest(t *testing.T) {
	manifest, err := scenario.LoadManifest(fixture("missing-manifest"))

	assert.NotNil(t, err)
	assert.Nil(t, manifest, "Returned a manifest")
	assert.Regexp(t, "scenarios.yaml: no such file or directory$", err.Error())
}

func Test_LoadManifest_manifest_missing_scenarios(t *testing.T) {
	manifest, err := scenario.LoadManifest(fixture("manifest-missing-scenarios"))

	assert.NotNil(t, err)
	assert.Nil(t, manifest, "Returned a manifest")
	assert.Regexp(t, "^Error unmarshalling", err.Error())
}
