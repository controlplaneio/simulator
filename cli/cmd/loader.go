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

func contains(scenarios []Scenario, name string) bool {
	for _, a := range scenarios {
		if a.Name == name {
			return true
		}
	}

	return false
}

func displayName(name string) string {
	return strings.Title(strings.Replace(name, "_", " ", -1))
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
			scenarios = append(scenarios, Scenario{Name: f.Name(), DisplayName: displayName(f.Name())})
		}
	}

	return scenarios, nil
}
