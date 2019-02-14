package commander

import (
	"github.com/gofunct/goexec/pkg/fs"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type CmdFunc func(c *Commander)

type Commander struct {
	root *cobra.Command
	debug *cobra.Command
	cron *cobra.Command
	script *cobra.Command
	cronz *cron.Cron
	*fs.Fs
}

func NewCommander(name, usg, version string) *Commander {
	f := fs.NewFs()
	cr := cron.New()

	cmder :=  &Commander{
		root: &cobra.Command{
			Use: name,
			Short: usg,
			Version: version,
		},
		Fs: f,
		debug: DebugCmd,
		script: ScriptCmd,
		cronz: cr,
		cron: CronCmd,
	}
	var dir string
	cmder.script.Flags().StringVar(&dir, "dir", os.Getenv("PWD"), "directory to execute script in")
}

func (c *Commander) AddDebugger(name, usg string, cmdFunc CmdFunc) {
	c.debug.AddCommand(&cobra.Command{
		Use: name,
		Short: usg,
		Run: func(cmd *cobra.Command, args []string) {
			cmdFunc(c)
		},
	})
}

func (c *Commander) AddScript(name, usg string, scrpt string) {
	exe := exec.Command("/bin/bash", "-c", scrpt)
	exe.Env = os.Environ()
	exe.Stderr = c.script.OutOrStderr()
	exe.Stdout = c.script.OutOrStdout()
	newscrpt := &cobra.Command{
		Use: name,
		Short: usg,
		Run: func(cmd *cobra.Command, args []string) {
			if err := exe.Run(); err != nil {
				panic(errors.Wrap(err, "failed to run script"))
			}
		},
	}
	new

	c.script.AddCommand(&cobra.Command{
		Use: name,
		Short: usg,
		Run: func(cmd *cobra.Command, args []string) {
			cmdFunc(c)
		},
	})
}

func (c *Commander) AddCron(name, usg string, cmdFunc CmdFunc) {

}