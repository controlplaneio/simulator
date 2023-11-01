package cli

import (
	"github.com/spf13/cobra"
)

var (
	bucket, key, name string
)

var simulatorCmd = &cobra.Command{
	Use: "simulator",
}

func Execute() {
	err := simulatorCmd.Execute()
	cobra.CheckErr(err)
}
