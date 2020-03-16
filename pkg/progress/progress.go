package progress

// TaskProgress represents a user's progress on a particular task
type TaskProgress struct {
	ID            int `json:"id"`
	LastHintIndex int `json:"lastHintIndex"`
	Score         int `json:"score"`
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
	return nil, nil
}

// SaveProrgess persists the user's progress on a scenairio to the local
// ~/.kubesim folder
func (lpp LocalStateProvider) SaveProgress(p ScenarioProgress) error {
	return nil
}
