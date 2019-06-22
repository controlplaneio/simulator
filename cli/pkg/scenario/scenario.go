package scenario

import (
	"fmt"
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
}

// Validate a scenario relative to its manifest
func (s *Scenario) Validate(manifestPath string) error {
	scenarioPath, err := filepath.Abs(filepath.Join(manifestPath, s.Path))
	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Error resolving %s from %s for scenario %s", s.Path, s.DisplayName, manifestPath))
	}

	stat, err := os.Stat(scenarioPath)
	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Error stating %s for scenario %s in %s", s.Path, s.DisplayName, manifestPath))
	}

	if stat.IsDir() != true {
		return errors.New(
			fmt.Sprintf("Scenario %s is not a directory at %s read from %s",
				s.DisplayName, s.Path, manifestPath))
	}

	return nil
}
