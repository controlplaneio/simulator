package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Scenario struct {
	Id          string `yaml:"id"`
	Path        string `yaml:"path"`
	DisplayName string `yaml:"name"`
}

type ScenarioManifest struct {
	Name      string     `yaml:"name"`
	Kind      string     `yaml:"kind"`
	Scenarios []Scenario `yaml:"scenarios"`
}

func contains(scenarios []Scenario, id string) bool {
	for _, a := range scenarios {
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

func manifestPath() string {
	var manifestPath = os.Getenv(manifestPathEnvVar)
	fmt.Println("Env for scenarios was " + manifestPath)
	if manifestPath == "" {
		manifestPath = defaultManifestPath
	}
	fmt.Println("Looking for scenarios in " + manifestPath)

	return manifestPath
}

func validateScenario(manifestPath string, scenario Scenario) error {
	scenarioPath, err := filepath.Abs(filepath.Join(manifestPath, scenario.Path))
	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Error resolving %s from %s for scenario %s", scenario.Path, scenario.DisplayName, manifestPath))
	}

	stat, err := os.Stat(scenarioPath)
	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Error stating %s for scenario %s in %s", scenario.Path, scenario.DisplayName, manifestPath))
	}

	if !stat.IsDir() {
		return errors.Wrap(err,
			fmt.Sprintf("Scenario %s is not a directory at %s read from %s",
				scenario.DisplayName, scenario.Path, manifestPath))
	}

	return nil
}

func loadScenarios(manifestPath string) ([]Scenario, error) {
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
	err = yaml.Unmarshal([]byte(manifestYaml), &manifest)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error unmarshalling %s", manifestPath))
	}

	for _, scenario := range manifest.Scenarios {
		err := validateScenario(manifestPath, scenario)
		if err != nil {
			return nil, err
		}
	}

	return manifest.Scenarios, nil
}
