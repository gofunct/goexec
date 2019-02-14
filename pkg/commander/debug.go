package commander

import "github.com/spf13/cobra"

var DebugCmd = &cobra.Command{
	Use: "debug",
	Short: "debug flags, config, etc",
}
