# scenario
--
    import "."

Package scenario is a package for loading scenario manifests from a
`scenarios.yaml` file and accessing and manipulating them programmatically.

## Usage

#### type Manifest

```go
type Manifest struct {
	// Name - the name of the manifest e.g. scenarios
	Name string `yaml:"name"`
	// Kind - unique name and version string idenitfying the schema of this
	// document
	Kind string `yaml:"kind"`
	// Scenarios - a list of Scenario structs representing the scenarios
	Scenarios []Scenario `yaml:"scenarios"`
}
```

Manifest structure representing a `scenarios.yaml` document

#### func  LoadManifest

```go
func LoadManifest(manifestPath string) (*Manifest, error)
```
LoadManifest loads a manifest named `scenarios.yaml` from the supplied path

#### func (*Manifest) Contains

```go
func (m *Manifest) Contains(id string) bool
```
Contains returns a boolean indicating whether a ScenarioManifest contains a
Scenario with the supplied id

#### func (*Manifest) Find

```go
func (m *Manifest) Find(id string) *Scenario
```
Find returns a scenario for the supplied id

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
