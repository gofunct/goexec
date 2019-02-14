package commander

import (
	"github.com/gofunct/goexec/pkg/commands"
	"github.com/gofunct/goexec/pkg/exec"
	"github.com/gofunct/goexec/pkg/fs"
	"github.com/gofunct/goexec/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

type Commander struct {
	root   *cobra.Command
	script *cobra.Command
	*fs.Fs
}

func NewCommander(name, usg, version string) *Commander {
	f := fs.NewFs()
	cmder := &Commander{
		root: &cobra.Command{
			Use:     name,
			Short:   usg,
			Version: version,
		},
		Fs: f,

		script: commands.ScriptCmd,
	}
	cmder.root.AddCommand(cmder.script, commands.DebugCmd, commands.LoadCmd)
	return cmder
}

func (c *Commander) AddScript(name, usg string, dir, script string) {
	script = c.Render(script)
	exe := exec.New()
	cmd := exe.Command(script)
	cmd.SetEnv(os.Environ())
	cmd.SetDir(dir)
	cmd.SetStdout(c.script.OutOrStdout())
	cmd.SetStderr(c.script.OutOrStderr())
	c.script.AddCommand(&cobra.Command{
		Use:   name,
		Short: usg,
		Long:  script,
		Run: func(_ *cobra.Command, args []string) {
			if err := cmd.Run(); err != nil {
				util.PrintErr(err, "failed to run script")
				os.Exit(1)
			}
		},
	})
}

func (c *Commander) AddDescription(s string) {
	c.root.Long = s
}

func (c *Commander) Execute() error {
	return c.root.Execute()
}
