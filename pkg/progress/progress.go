package progress

type TaskProgress struct {
	Score int `json:"score"`
}

type ScenarioProgress struct {
	Name        string         `json:"name"`
	CurrentTask int            `json:"currentTask"`
	Tasks       []TaskProgress `json:"tasks"`
}

type Progress struct {
	Scenarios []ScenarioProgress `json:"scenarioProgress"`
}
