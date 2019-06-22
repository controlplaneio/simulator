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

var badManifestTests = []struct {
	name         string
	errorPattern string
}{
	{"missing-manifest", "scenarios.yaml: no such file or directory$"},
	{"manifest-missing-scenarios", "^Error unmarshalling"},
	{"malformed-manifest", "^Error unmarshalling"},
}

func Test_LoadManifest_bad_manifests(t *testing.T) {
	for _, tt := range badManifestTests {
		t.Run(tt.name, func(t *testing.T) {
			manifest, err := scenario.LoadManifest(fixture(tt.name))

			assert.NotNil(t, err)
			assert.Nil(t, manifest, "Returned a manifest")
			assert.Regexp(t, tt.errorPattern, err.Error())
		})
	}
}
