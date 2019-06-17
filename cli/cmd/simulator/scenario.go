package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const scenariosPath = "./simulation-scripts/scenario/"

type Scenario struct {
	Name        string
	DisplayName string
}

func getScenarios() ([]Scenario, error) {
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

func init() {
	rootCmd.AddCommand(scenarioCmd)
}

var scenarioCmd = &cobra.Command{
	Use:   `scenario`,
	Short: "Lists available scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		scenarios, err := getScenarios()

		if err != nil {
			return err
		}

		for _, s := range scenarios {
			fmt.Println(s.DisplayName)
		}

		return nil
	},
}
