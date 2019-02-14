package commander

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	dir string
)

func init() {
	ScriptCmd.PersistentFlags().StringVarP(&dir, "dir", "d", os.Getenv("PWD"), "directory to run script in")
}

var ScriptCmd = &cobra.Command{
	Use: "script",
	Short: "run a bash script",
}