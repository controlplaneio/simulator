package progress

import (
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
type LocalStateProvider struct {
	Logger *logrus.Logger
}

// NewLocalStateProvider returns an instance of LocalStateProvider
func NewLocalStateProvider(logger *logrus.Logger) LocalStateProvider {
	return LocalStateProvider{
		Logger: logger,
	}
}

// StateProvider defines the contract for retrieving and persisting a user's
// progress
type StateProvider interface {
	GetLogger() *logrus.Logger
	GetProgress(scenario string) (*ScenarioProgress, error)
	SaveProgress(p ScenarioProgress) error
}

// fileExists checks that a filename exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		// stat syscall failed, this is bad
		panic(err)
	}

	return !info.IsDir()
}

// GetLogger returns the state provider's logger
func (lsp LocalStateProvider) GetLogger() *logrus.Logger {
	return lsp.Logger
}

func (lsp LocalStateProvider) writeProgress(p *Progress) error {
	path, err := util.ExpandTilde(ProgressPath)
	if err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error resolving progress path")
		return errors.Wrap(err, "Error resolving progress path")
	}

	data, err := json.Marshal(&p)
	if err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error":    err,
			"Progress": p,
		}).Fatal("Error marshaling progres to JSON")
		panic(err)
	}

	if err = ioutil.WriteFile(*path, data, 0660); err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error":    err,
			"Path":     *path,
			"Progress": p,
		}).Error("Error writing progress to disk")
		return errors.Wrap(err, "Error writing progress file")
	}

	return nil
}

func (lsp LocalStateProvider) getProgress() (*Progress, error) {
	path, err := util.ExpandTilde(ProgressPath)
	if err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error resolving progress path")
		return nil, errors.Wrap(err, "Error resolving progress path")
	}

	if !fileExists(*path) {
		p := Progress{}
		if err = lsp.writeProgress(&p); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(*path)
	if err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error opening progress file")
		return nil, errors.Wrap(err, "Error opening progress file")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error reading progress file")
		return nil, errors.Wrap(err, "Error reading progress file")
	}
	var p Progress

	if err = json.Unmarshal(bytes, &p); err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error unmarshaling progress JSON")
		return nil, errors.Wrap(err, "Error unmarshaling progress JSON")
	}

	return &p, nil
}

// GetProgress retrieves the user's progress for the provided scenario
func (lsp LocalStateProvider) GetProgress(scenario string) (*ScenarioProgress, error) {
	p, err := lsp.getProgress()
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
		lsp.Logger.WithFields(logrus.Fields{
			"Scenario": scenario,
			"Progress": retVal,
		}).Info("Found existing progress")
		return &retVal, nil
	}

	lsp.Logger.WithFields(logrus.Fields{
		"Scenario": scenario,
		"Progress": retVal,
	}).Info("No existing progress found")

	return nil, nil
}

func remove(slice []ScenarioProgress, i int) []ScenarioProgress {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// SaveProgress persists the user's progress on a scenairio to the local
// ~/.kubesim folder
func (lsp LocalStateProvider) SaveProgress(update ScenarioProgress) error {
	p, err := lsp.getProgress()
	if err != nil {
		return err
	}

	found := false
	var existing ScenarioProgress
	index := 0
	for i, sp := range p.Scenarios {
		if sp.Name == update.Name {
			found = true
			index = i
			break
		}
	}

	if found {
		lsp.Logger.WithFields(logrus.Fields{
			"Scenario":         update.Name,
			"ExistingProgress": existing,
		}).Info("Found existing progress to remove")
		p.Scenarios = remove(p.Scenarios, index)
	}

	lsp.Logger.WithFields(logrus.Fields{
		"Scenario":    update.Name,
		"NewProgress": update,
	}).Info("Adding new progress")
	p.Scenarios = append(p.Scenarios, update)

	if err = lsp.writeProgress(p); err != nil {
		lsp.Logger.WithFields(logrus.Fields{
			"Scenario":    update.Name,
			"NewProgress": update,
			"Error":       err,
		}).Error("Error writing progress")
		return err
	}

	return nil
}
