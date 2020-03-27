package progress_test

import (
	"github.com/go-test/deep"
	"github.com/kubernetes-simulator/simulator/pkg/progress"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/stretchr/testify/assert"

	"bytes"
	"os"
	"testing"
	"text/template"
)

var progressTemplateSrc = `{
	"scenarioProgress": [
		{
			"Name": "{{.Name}}",
			"currentTask": 1,
			"tasks": [
				{ "id": 1, "lastHintIndex": 2, "score": null }
			]
		}
	]
}`

func makeProgress(name string) {
	_ = os.Remove(util.MustExpandTilde(progress.ProgressPath))
	progressTmpl, err := template.New("progress-template").Parse(progressTemplateSrc)
	if err != nil {
		panic(err)

	}

	var buf bytes.Buffer

	err = progressTmpl.Execute(&buf, struct{ Name string }{Name: name})
	if err != nil {
		panic(err)
	}

	if _, err = util.EnsureFile(util.MustExpandTilde(progress.ProgressPath), buf.String()); err != nil {
		panic(err)
	}
}

func Test_GetProgress_with_existing_progress(t *testing.T) {
	lsp := progress.NewLocalStateProvider(NullLogger())
	name := "test-scenario"
	makeProgress(name)

	actual, err := lsp.GetProgress(name)
	if err != nil {
		t.Fatal(err)
	}

	var currentTask, lastHintIndex *int
	currentTask = new(int)
	*currentTask = 1

	lastHintIndex = new(int)
	*lastHintIndex = 2

	assert.Equal(t, progress.ScenarioProgress{
		Name:        name,
		CurrentTask: currentTask,
		Tasks: []progress.TaskProgress{
			progress.TaskProgress{ID: 1, LastHintIndex: lastHintIndex, Score: nil},
		},
	}, *actual, "Expected matching scenario progress to be returned")

}

func Test_GetProgress_with_no_existing_progress(t *testing.T) {
	_ = os.Remove(util.MustExpandTilde(progress.ProgressPath))
	lsp := progress.NewLocalStateProvider(NullLogger())

	actual, err := lsp.GetProgress("test-scenario")
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err, "Expected no error")
	assert.Nil(t, actual, "Expected no progress")
}

func makeScenarioProgress(name string) progress.ScenarioProgress {
	score := 100
	var currentTask, lastHintIndex *int
	currentTask = new(int)
	*currentTask = 2

	lastHintIndex = new(int)
	*lastHintIndex = 1

	sp := progress.ScenarioProgress{
		Name:        name,
		CurrentTask: currentTask,
		Tasks: []progress.TaskProgress{
			progress.TaskProgress{
				ID:             1,
				LastHintIndex:  lastHintIndex,
				Score:          &score,
				ScoringSkipped: false},
		},
	}

	return sp
}

func Test_SaveProgress_with_no_existing_Progress(t *testing.T) {
	_ = os.Remove(util.MustExpandTilde(progress.ProgressPath))
	lsp := progress.NewLocalStateProvider(NullLogger())
	name := "test-scenario"
	sp := makeScenarioProgress(name)
	err := lsp.SaveProgress(sp)
	assert.Nil(t, err, "Expected no error saving progress")

	actual, err := lsp.GetProgress(sp.Name)
	assert.Nil(t, err, "Expected no error getting progress")
	if diff := deep.Equal(sp, *actual); diff != nil {
		t.Error("Returned progress did not match saved progress", diff)
	}

}

func Test_SaveProgress_with_existing_Progress(t *testing.T) {
	_ = os.Remove(util.MustExpandTilde(progress.ProgressPath))
	lsp := progress.NewLocalStateProvider(NullLogger())
	name := "test-scenario"

	makeProgress(name)
	sp := makeScenarioProgress(name)

	err := lsp.SaveProgress(sp)
	assert.Nil(t, err, "Expected no error saving progress")

	actual, err := lsp.GetProgress(sp.Name)
	assert.Nil(t, err, "Expected no error getting progress")
	if diff := deep.Equal(sp, *actual); diff != nil {
		t.Error("Returned progress did not match saved progress", diff)
	}
}
