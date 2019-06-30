package scenario

import (
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Manifest structure representing a `scenarios.yaml` document
type Manifest struct {
	// Name - the name of the manifest e.g. scenarios
	Name string `yaml:"name"`
	// Kind - unique name and version string idenitfying the schema of this document
	Kind string `yaml:"kind"`
	// Scenarios - a list of Scenario structs representing the scenarios
	Scenarios []Scenario `yaml:"scenarios"`
}

// Contains returns a boolean indicating whether a ScenarioManifest contains a Scenario
// with the supplied id
func (m *Manifest) Contains(id string) bool {
	for _, a := range m.Scenarios {
		if a.Id == id {
			return true
		}
	}

	return false
}

// Find returns a scenario for the supplied id
func (m *Manifest) Find(id string) *Scenario {
	for _, a := range m.Scenarios {
		if a.Id == id {
			return &a
		}
	}

	return nil
}

const (
	manifestPathEnvVar  = "SIMULATOR_MANIFEST_PATH"
	defaultManifestPath = "./simulation-scripts/"
	manifestFileName    = "scenarios.yaml"
)

// ManifestPath reads the manifest path from the environment variable `SIMULATOR_MANIFEST_PATH`
// or uses a default value of `../simulation-scripts`
func ManifestPath() string {
	var manifestPath = os.Getenv(manifestPathEnvVar)
	if manifestPath == "" {
		manifestPath = defaultManifestPath
	}

	return manifestPath
}

// LoadManifest loads a manifest named `scenarios.yaml` from the supplied path
func LoadManifest(manifestPath string) (*Manifest, error) {
	joined := filepath.Join(manifestPath, manifestFileName)
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Error resolving manifest file %s from %s", manifestFileName, manifestPath)
	}

	manifestYaml, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading manifest file %s", manifestPath)
	}

	manifest := Manifest{}
	err = yaml.UnmarshalStrict([]byte(manifestYaml), &manifest)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshalling %s", manifestPath)
	}

	if structs.HasZero(manifest) {
		return nil, errors.Errorf("Error unmarshalling %s - missing required fields", manifestPath)
	}

	for _, scenario := range manifest.Scenarios {
		err = scenario.Validate(manifestPath)
		if err != nil {
			return nil, err
		}
	}

	return &manifest, nil
}
