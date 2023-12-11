package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionInfo contains the version information.
type VersionInfo struct {
	Version   string // the semantic version, grabbed from git tags during build
	AppName   string // name of the application
	GitHash   string // the git hash of the build
	BuildDate string // build date, will be injected by the build system
}

// WithVersionCmd creates a new version command
func WithVersionCmd(ver VersionInfo) SimulatorCmdOptions {
	return func(simulator *cobra.Command) {
		versionCmd := &cobra.Command{
			Use:   "version",
			Short: "Display the version information",
			Run: func(cmd *cobra.Command, args []string) {
				//nolint:forbidigo
				fmt.Printf("%s version: %s\n", ver.AppName, ver.Version)
				//nolint:forbidigo
				fmt.Printf("Git commit hash: %s\n", ver.GitHash)
				//nolint:forbidigo
				fmt.Printf("Build date: %s\n", ver.BuildDate)
			},
		}
		simulator.AddCommand(versionCmd)
	}
}
