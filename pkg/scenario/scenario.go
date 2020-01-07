package scenario

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// Scenario structure representing a scenario
type Scenario struct {
	// A machine parseable unique id for the scenario
	Id string `yaml:"id"`
	// Path to the scenario - paths are relative to the ScenarioManifest that
	// defines this scenario
	Path string `yaml:"path"`
	// A human-friendly readable name for this scenario for use in user interfaces
	DisplayName string `yaml:"name"`
	// A Difficulty level for this scenario for use in user interfaces
	Difficulty string `yaml:"difficulty"`
	// A short description of the scenario to be used in the user interfaces
	Description string `yaml:"description"`
}

// Validate a scenario relative to its manifest
func (s *Scenario) Validate(manifestPath string) error {
	scenarioPath, err := filepath.Abs(filepath.Join(manifestPath, s.Path))
	if err != nil {
		return errors.Wrapf(err,
			"Error resolving %s from %s for scenario %s", s.Path, s.DisplayName, manifestPath)
	}

	stat, err := os.Stat(scenarioPath)
	if err != nil {
		return errors.Wrapf(err,
			"Error stating %s for scenario %s in %s", s.Path, s.DisplayName, manifestPath)
	}

	if !stat.IsDir() {
		return errors.Errorf(
			"Scenario %s is not a directory at %s read from %s",
			s.DisplayName, s.Path, manifestPath)
	}

	return nil
}
