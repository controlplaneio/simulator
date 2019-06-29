package scenario_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/scenario"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LoadManifest(t *testing.T) {
	t.Parallel()
	manifest, err := scenario.LoadManifest(fixture("valid"))

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

func Test_Contains(t *testing.T) {
	t.Parallel()
	m := scenario.Manifest{
		Name: "test",
		Kind: "test/0.1",
		Scenarios: []scenario.Scenario{
			scenario.Scenario{
				Id:          "test_scenario",
				DisplayName: "Test Scenario",
				Path:        "./test",
			},
		},
	}

	assert.True(t, m.Contains("test_scenario"), "Contains did not return true for valid scenario")
	assert.False(t, m.Contains("invalid"), "Contains did not return false for invalid scenario")
}

func Test_Find(t *testing.T) {
	t.Parallel()
	s := scenario.Scenario{
		Id:          "test_scenario",
		DisplayName: "Test Scenario",
		Path:        "./test",
	}
	m := scenario.Manifest{
		Name:      "test",
		Kind:      "test/0.1",
		Scenarios: []scenario.Scenario{s},
	}

	assert.Equal(t, m.Find("test_scenario"), &s, "Contains did not return valid scenario")
	assert.Nil(t, m.Find("invalid"), "Contains returned scenario for invalid scenario")
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
	{"missing-scenario", "^Error stating"},
	{"scenario-not-a-dir", "is not a directory"},
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
