package cmd

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Scenario struct {
	Id          string
	DirName     string
	DisplayName string
}

func contains(scenarios []Scenario, id string) bool {
	for _, a := range scenarios {
		if a.Id == id {
			return true
		}
	}

	return false
}

func displayName(name string) string {
	return strings.Title(strings.Replace(name, "_", " ", -1))
}

func makeScenario(dirName string) Scenario {
	return Scenario{
		DirName:     dirName,
		DisplayName: displayName(dirName),
		Id:          strings.ToLower(dirName),
	}
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
			scenarios = append(scenarios, makeScenario(f.Name()))
		}
	}

	return scenarios, nil
}
