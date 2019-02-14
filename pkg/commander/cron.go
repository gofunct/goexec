package commander

import "github.com/spf13/cobra"

var CronCmd = &cobra.Command{
	Use: "short",
	Short: "run cron jobs on specified schedule",
}
