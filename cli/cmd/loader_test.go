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

func Test_makeScenario(t *testing.T) {
	dirName := "lazy_AWS"
	s := makeScenario(dirName)

	if s.DirName != dirName {
		t.Errorf("Did not set DirName to supplied argument")
	}

	if s.DisplayName != "Lazy AWS" {
		t.Errorf("Did not set DisplayName correctly")
	}

	if s.Id != "lazy_aws" {
		t.Errorf("Did not set ID correctly")
	}
}
