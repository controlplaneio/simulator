package cmd

import (
	"testing"
)

func Test_loadScenarios(t *testing.T) {
	scenarios, err := loadScenarios("../../simulation-scripts/scenario")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(scenarios) == 0 {
		t.Errorf("Returned no scenarios")
	}
}
