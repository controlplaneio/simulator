package scenarios

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v2"

	"github.com/controlplaneio/simulator/v2/controlplane"
)

//go:embed scenarios.yaml
var config embed.FS

func List() ([]Scenario, error) {
	var scenarios []Scenario

	b, err := config.ReadFile("scenarios.yaml")
	if err != nil {
		slog.Error("failed to load scenarios file")
		return nil, err
	}

	err = yaml.Unmarshal(b, &scenarios)
	if err != nil {
		slog.Error("failed to unmarshall scenarios")
		return nil, err
	}

	return scenarios, nil
}

func Find(id string) (Scenario, error) {
	var s Scenario

	scenarios, err := List()
	if err != nil {
		return s, err
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

func (s Scenario) Challenge() ([]byte, error) {
	return config.ReadFile(fmt.Sprintf("roles/%s/files/challenge.txt", s.ID))
}

func (s Scenario) Playbook() string {
	return fmt.Sprintf("%s/%s", controlplane.AnsiblePlaybookDir, s.ID)
}
