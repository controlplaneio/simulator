package scenarios

import (
	"embed"
	"errors"
	"log/slog"

	"gopkg.in/yaml.v2"
)

//go:embed scenarios.yaml
var config embed.FS

func List() ([]Scenario, error) {
	var scenarios []Scenario

	b, err := config.ReadFile("scenarios.yaml")
	if err != nil {
		slog.Error("failed to load scenarios file")
		return nil, errors.Join(errors.New("failed to list scenarios"), err)
	}

	err = yaml.Unmarshal(b, &scenarios)
	if err != nil {
		slog.Error("failed to unmarshall scenarios")
		return nil, errors.Join(errors.New("failed to list scenarios"), err)
	}

	return scenarios, nil
}

func Find(id string) (Scenario, error) {
	var s Scenario

	scenarios, err := List()
	if err != nil {
		return s, errors.Join(errors.New("failed to find scenario"), err)
	}

	for _, scenario := range scenarios {
		if scenario.ID == id {
			return scenario, nil
		}
	}

	return s, errors.New("unable to find scenario")
}

type Scenario struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Difficulty  string `yaml:"difficulty"`
}
