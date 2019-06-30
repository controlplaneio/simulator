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
EnvOrDefault tries to read the key and returns a default value if it is empty

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

#### func  Slurp

```go
func Slurp(path string) (*string, error)
```
Slurp reads an entire file into a string in one operation and returns a pointer
to the file's content
