package progress_test

import (
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
	_ = os.Remove(progress.ProgressPath)
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

func Test_GetProgress(t *testing.T) {
	lsp := progress.LocalStateProvider{}
	name := "test-scenario"
	makeProgress(name)

	actual, err := lsp.GetProgress(name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, progress.ScenarioProgress{
		Name:        name,
		CurrentTask: 1,
		Tasks: []progress.TaskProgress{
			progress.TaskProgress{ID: 1, LastHintIndex: 2, Score: nil},
		},
	}, *actual, "Expected matching scenario progress to be returned")

}
