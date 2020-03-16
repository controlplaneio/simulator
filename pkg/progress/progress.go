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
	LastHintIndex  int  `json:"lastHintIndex"`
	Score          *int `json:"score"`
	ScoringSkipped bool `json:"scoringSkipped"`
}

// ScenarioProgress represents a user's progress on a particular scenario
type ScenarioProgress struct {
	Name        string         `json:"name"`
	CurrentTask int            `json:"currentTask"`
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

// GetProgress retrieves the user's progress for the provided scenario
func (lpp LocalStateProvider) GetProgress(scenario string) (*ScenarioProgress, error) {
	path, err := util.ExpandTilde(ProgressPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error resolving progress path")
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

// SaveProgress persists the user's progress on a scenairio to the local
// ~/.kubesim folder
func (lpp LocalStateProvider) SaveProgress(p ScenarioProgress) error {
	return nil
}
