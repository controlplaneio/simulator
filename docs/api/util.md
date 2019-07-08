# util
--
    import "."


## Usage

#### func  Debug

```go
func Debug(msg ...interface{})
```
Debug prints a debug message to stdout

#### func  DetectPublicIP

```go
func DetectPublicIP() (*string, error)
```
DetectPublicIP detects your public IP address

#### func  EnsureFile

```go
func EnsureFile(path, contents string) (bool, error)
```
EnsureFile checks a file exists and writes the supplied contents if not. returns
a boolean indicating whether it wrote a file or not and any error

#### func  EnvOrDefault

```go
func EnvOrDefault(key, def string) string
```
EnvOrDefault tries to read an environment variable with the supplied key and
returns its value. EnvOrDefault returns a default value if it is empty or unset

#### func  ExpandTilde

```go
func ExpandTilde(path string) (*string, error)
```
ExpandTilde returns the fully qualified path to a file in the user's home
directory. I.E. it expands a path beginning with `~/`) and checks the file
exists. ExpandTilde will cache the user's home directory to amortise the cost of
the syscall

#### func  FileExists

```go
func FileExists(path string) (bool, error)
```
FileExists checks whether a path exists

#### func  MustRemove

```go
func MustRemove(path string)
```
MustRemove removes a file or empty directory. MustRemove will ignore an error if
the path doesn't exist or panic for any other error

#### func  MustSlurp

```go
func MustSlurp(path string) string
```
MustSlurp is the panicky counterpart to Slurp. MustSlurp reads an entire file
into a string in one operation and returns the contents or panics if it
encouters and error

#### func  Run

```go
func Run(wd string, env []string, cmd string, args ...string) (*string, error)
```
Run runs a child process and returns its buffer stdout. Run also tees the output
to stdout of this process, `env` will be appended to the current environment.
`wd` is the working directory for the child

#### func  RunSilently

```go
func RunSilently(wd string, env []string, cmd string, args ...string) (*string, *string, error)
```
RunSilently runs a sub command silently

#### func  SSH

```go
func SSH(host string) error
```

#### func  Slurp

```go
func Slurp(path string) (*string, error)
```
Slurp reads an entire file into a string in one operation and returns a pointer
to the file's content or an error if any. Similar to `ioutil.ReadFile` but it
calls `filepath.Abs` first which cleans the path and resolves relative paths
from the working directory.

Note that this is slightly less efficient for zero-length files than
`ioutil.Readfile` as it uses the default read buffer size of `bytes.MinRead`
internally

#### func  StartInteractiveShell

```go
func StartInteractiveShell(sshConfig *ssh.ClientConfig, network string, host string, port string) error
```
