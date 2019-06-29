# runner
--
    import "."


## Usage

#### func  Create

```go
func Create() error
```
Create runs terraform init, plan, apply to create the necessary infratsructure
to run scenarios

#### func  Destroy

```go
func Destroy() error
```
Destroy call terraform destroy to remove the infrastructure

#### func  EnsureTfVarsFile

```go
func EnsureTfVarsFile(tfDir string, publicKey string, accessCIDR string) error
```
EnsureTfVarsFile writes an tfvars file if one hasnt already been made

#### func  EnvOrDefault

```go
func EnvOrDefault(key string, def string) string
```
EnvOrDefault tries to read the key and returns a default value if it is empty

#### func  FileExists

```go
func FileExists(path string) (bool, error)
```
FileExists checks whether a path exists

#### func  InitIfNeeded

```go
func InitIfNeeded() error
```
InitIfNeeded checks if there is a terraform state folder and calls terraform
init if not

#### func  Perturb

```go
func Perturb(po *PerturbOptions) (*string, error)
```
Perturb runs the perturb script with the supplied options

#### func  PrepareTfArgs

```go
func PrepareTfArgs(cmd string) []string
```
PrepareTfArgs takes a string with the terraform command desired and returns a
slice of strings containing the complete list of arguments including the command
to use when exec'ing terraform

#### func  ReadFile

```go
func ReadFile(path string) (*string, error)
```
ReadFile return a pointer to a string with the file's content

#### func  Run

```go
func Run(wd string, env []string, cmd string, args ...string) (*string, error)
```
Run runs a child process and returns its buffer stdout. Run also tees the output
to stdout of this process, `env` will be appended to the current environment.
`wd` is the working directory for the child

#### func  Terraform

```go
func Terraform(cmd string) (*string, error)
```
Terraform wraps running terraform as a child process

#### type PerturbOptions

```go
type PerturbOptions struct {
	Master       net.IP
	Slaves       []net.IP
	ScenarioName string
}
```

PerturbOptions represents the parameters required by the perturb.sh script

#### func  MakePerturbOptions

```go
func MakePerturbOptions(tfo TerraformOutput, path string) PerturbOptions
```
MakePerturbOptions takes a TerraformOutput and a path to a scenario and makes a
struct of PerturbOptions

#### func (*PerturbOptions) String

```go
func (po *PerturbOptions) String() string
```

#### func (*PerturbOptions) ToArguments

```go
func (po *PerturbOptions) ToArguments() []string
```
ToArguments converts a PerturbOptions struct into a slice of strings containing
the command line options to pass to perturb

#### type SSHConfig

```go
type SSHConfig struct {
	Hostname    string
	KeyFilePath string
	User        string
	BastionIP   string
}
```

SSHConfig represents the values needed to produce a config block to allow SSH to
the private kubernetes nodes via the bastion

#### type StringOutput

```go
type StringOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}
```

StringOutput is a struct representing an output from terraform that contains a
string

#### type StringSliceOutput

```go
type StringSliceOutput struct {
	Sensitive bool          `json:"sensitive"`
	Type      []interface{} `json:"type"`
	Value     []string      `json:"value"`
}
```

StringSliceOutput is a struct representing an output from terraform that
contains a slice of strings

#### type TerraformOutput

```go
type TerraformOutput struct {
	BastionPublicIP       StringOutput      `json:"bastion_public_ip"`
	ClusterNodesPrivateIP StringSliceOutput `json:"cluster_nodes_private_ip"`
	MasterNodesPrivateIP  StringSliceOutput `json:"master_nodes_private_ip"`
}
```

TerraformOutput is a struct representing the expected output variables from the
terraform script

#### func  ParseTerraformOutput

```go
func ParseTerraformOutput(output string) (*TerraformOutput, error)
```
ParseTerraformOutput takes a string containing the stdout from `terraform output
-json` and returns a TerraformOutput struct

#### func  Status

```go
func Status() (*TerraformOutput, error)
```
Status calls terraform output to get the state of the infrastruture and parses
the output for programmatic use

#### func (*TerraformOutput) IsUsable

```go
func (tfo *TerraformOutput) IsUsable() bool
```
IsUsable checks whether the TerraformOutput has all the necessary information to
be converted for use with perturb

#### func (*TerraformOutput) ToSSHConfig

```go
func (tfo *TerraformOutput) ToSSHConfig() (*string, error)
```
ToSSHConfig produces the SSH config

#### type TfVars

```go
type TfVars struct {
	PublicKey  string
	AccessCIDR string
}
```

TfVars struct representing the input variables for terraform to create the
infrastructure

#### func  NewTfVars

```go
func NewTfVars(publicKey string, accessCIDR string) TfVars
```
NewTfVars creates a TfVars struct with all the defaults

#### func (*TfVars) String

```go
func (tfv *TfVars) String() string
```
