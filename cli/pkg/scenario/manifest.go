package scenario

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ScenarioManifest structure representing a `scenarios.yaml` document
type ScenarioManifest struct {
	// Name - the name of the manifest e.g. scenarios
	Name string `yaml:"name"`
	// Kind - unique name and version string idenitfying the schema of this document
	Kind string `yaml:"kind"`
	// Scenarios - a list of Scenario structs representing the scenarios
	Scenarios []Scenario `yaml:"scenarios"`
}

// Returns a boolean indicating whether a ScenarioManifest contains a Scenario
// with the supplied id
func (m *ScenarioManifest) Contains(id string) bool {
	for _, a := range m.Scenarios {
		if a.Id == id {
			return true
		}
	}

	return false
}

const (
	manifestPathEnvVar  = "SIMULATOR_MANIFEST_PATH"
	defaultManifestPath = "../simulation-scripts/"
	manifestFileName    = "scenarios.yaml"
)

// Reads the manifest path from the environment variable `SIMULATOR_MANIFEST_PATH`
// or uses a default value of `../simulation-scripts`
func ManifestPath() string {
	var manifestPath = os.Getenv(manifestPathEnvVar)
	if manifestPath == "" {
		manifestPath = defaultManifestPath
	}

	return manifestPath
}

// Loads a manifest named scenarios.yaml from the supplied path
func LoadManifest(manifestPath string) (*ScenarioManifest, error) {
	joined := filepath.Join(manifestPath, manifestFileName)
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("Error resolving manifest file %s from %s", manifestFileName, manifestPath))
	}

	manifestYaml, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error reading manifest file %s", manifestPath))
	}

	manifest := ScenarioManifest{}
	err = yaml.UnmarshalStrict([]byte(manifestYaml), &manifest)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error unmarshalling %s", manifestPath))
	}

	if structs.HasZero(manifest) {
		return nil, errors.New(fmt.Sprintf("Error unmarshalling %s - missing required fields", manifestPath))
	}

	for _, scenario := range manifest.Scenarios {
		err = scenario.Validate(manifestPath)
		if err != nil {
			return nil, err
		}
	}

	return &manifest, nil
}
