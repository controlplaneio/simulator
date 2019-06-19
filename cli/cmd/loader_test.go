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

func Test_displayName(t *testing.T) {
	simpleName := displayName("foo")
	if simpleName != "Foo" {
		t.Errorf("Did not make simple name title case")
	}

	compoundName := displayName("hack_the_world")
	if compoundName != "Hack The World" {
		t.Errorf("Did not make compound name title case replacing the underscores")
	}
}
