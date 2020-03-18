package progress

import (
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"

	"encoding/json"
	"io/ioutil"
	"os"
)

// ProgressPath is the path to the progress file
const ProgressPath = "~/.kubesim/progress.json"

// TaskProgress represents a user's progress on a particular task
type TaskProgress struct {
	ID             int  `json:"id"`
	LastHintIndex  *int `json:"lastHintIndex"`
	Score          *int `json:"score"`
	ScoringSkipped bool `json:"scoringSkipped"`
}

// ScenarioProgress represents a user's progress on a particular scenario
type ScenarioProgress struct {
	Name        string         `json:"name"`
	CurrentTask *int           `json:"currentTask"`
	Tasks       []TaskProgress `json:"tasks"`
}

// Progress represents a user's progress for all scenarios they have launched
type Progress struct {
	Scenarios []ScenarioProgress `json:"scenarioProgress"`
}

// LocalStateProvider persists and retrieves a user's progress to the local
// ~/.kubesim folder
type LocalStateProvider struct{}

// StateProvider defines the contract for retrieving and persisting a user's
// progress
type StateProvider interface {
	GetProgress(scenario string) (*ScenarioProgress, error)
	SaveProgress(p ScenarioProgress) error
}

// fileExists checks that a filename exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func writeProgress(p *Progress) error {
	path, err := util.ExpandTilde(ProgressPath)
	if err != nil {
		return errors.Wrap(err, "Error resolving progress path")
	}

	data, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(*path, data, 0660); err != nil {
		return errors.Wrap(err, "Error writing progress file")
	}

	return nil
}

func getProgress() (*Progress, error) {
	path, err := util.ExpandTilde(ProgressPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error resolving progress path")
	}

	if !fileExists(*path) {
		p := Progress{}
		if err = writeProgress(&p); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(*path)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening progress file")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading progress file")
	}
	var p Progress

	if err = json.Unmarshal(bytes, &p); err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling progress json")
	}

	return &p, nil
}

// GetProgress retrieves the user's progress for the provided scenario
func (lsp LocalStateProvider) GetProgress(scenario string) (*ScenarioProgress, error) {
	p, err := getProgress()
	if err != nil {
		return nil, err
	}

	var retVal ScenarioProgress
	found := false
	for _, sp := range p.Scenarios {
		if sp.Name == scenario {
			retVal = sp
			found = true
			break
		}
	}

	if found {
		return &retVal, nil
	}

	return nil, nil
}

func remove(slice []ScenarioProgress, i int) []ScenarioProgress {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// SaveProgress persists the user's progress on a scenairio to the local
// ~/.kubesim folder
func (lsp LocalStateProvider) SaveProgress(update ScenarioProgress) error {
	p, err := getProgress()
	if err != nil {
		return err
	}

	found := false
	index := 0
	for i, sp := range p.Scenarios {
		if sp.Name == update.Name {
			found = true
			index = i
			break
		}
	}

	if found {
		p.Scenarios = remove(p.Scenarios, index)
	}

	p.Scenarios = append(p.Scenarios, update)

	if err = writeProgress(p); err != nil {
		return err
	}

	return nil
}
