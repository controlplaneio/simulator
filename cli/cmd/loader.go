package cmd

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Scenario struct {
	Name        string
	DisplayName string
}

func loadScenarios(scenariosPath string) ([]Scenario, error) {
	absPath, err := filepath.Abs(scenariosPath)

	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(absPath)

	if err != nil {
		return nil, err
	}

	var scenarios []Scenario
	for _, f := range files {
		if f.IsDir() {
			var displayName = strings.Title(strings.Replace(f.Name(), "_", " ", -1))
			scenarios = append(scenarios, Scenario{Name: f.Name(), DisplayName: displayName})
		}
	}

	return scenarios, nil
}
