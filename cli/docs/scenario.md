# scenario
--
    import "."


## Usage

#### func  ManifestPath

```go
func ManifestPath() string
```
Reads the manifest path from the environment variable `SIMULATOR_MANIFEST_PATH`
or uses a default value of `../simulation-scripts`

#### type Scenario

```go
type Scenario struct {
	// A machine parseable unique id for the scenario
	Id string `yaml:"id"`
	// Path to the scenario - paths are relative to the ScenarioManifest that
	// defines this scenario
	Path string `yaml:"path"`
	// A human-friendly readable name for this scenario for use in user interfaces
	DisplayName string `yaml:"name"`
}
```

Scenario structure representing a scenario

#### func (*Scenario) Validate

```go
func (s *Scenario) Validate(manifestPath string) error
```
Validate a scenario relative to its manifest

#### type ScenarioManifest

```go
type ScenarioManifest struct {
	// Name - the name of the manifest e.g. scenarios
	Name string `yaml:"name"`
	// Kind - unique name and version string idenitfying the schema of this document
	Kind string `yaml:"kind"`
	// Scenarios - a list of Scenario structs representing the scenarios
	Scenarios []Scenario `yaml:"scenarios"`
}
```

ScenarioManifest structure representing a `scenarios.yaml` document

#### func  LoadManifest

```go
func LoadManifest(manifestPath string) (*ScenarioManifest, error)
```
Loads a manifest named scenarios.yaml from the supplied path

#### func (*ScenarioManifest) Contains

```go
func (m *ScenarioManifest) Contains(id string) bool
```
Returns a boolean indicating whether a ScenarioManifest contains a Scenario with
the supplied id
