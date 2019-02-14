package commander

import (
	"fmt"
	sv2 "github.com/Masterminds/semver"
	"github.com/gofunct/goexec/pkg/commands"
	"github.com/gofunct/goexec/pkg/exec"
	"github.com/gofunct/goexec/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

type Commander struct {
	root   *cobra.Command
	script *cobra.Command
}

func NewCommander(name, usg string) *Commander {
	cmder := &Commander{
		root: &cobra.Command{
			Use:   util.BlueString("%s", name),
			Short: util.BlueString("%s", usg),
		},
		script: commands.ScriptCmd,
	}
	cmder.root.SetOutput(util.ColoredOut)
	cmder.script.SetOutput(cmder.root.OutOrStderr())
	commands.LoadCmd.SetOutput(cmder.root.OutOrStderr())
	commands.DebugCmd.SetOutput(cmder.root.OutOrStderr())
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
		Use:   util.BlueString("%s", name),
		Short: util.BlueString("%s", usg),
		Long:  util.BlueString("%s", "Script: "+script),
		Run: func(_ *cobra.Command, args []string) {
			util.Panic(cmd.Run(), "failed to run script: %v\n", script)
		},
	})
}

func (c *Commander) AddDescription(s string) {
	c.root.Long = s
}

func (c *Commander) AddVersion(s string) {
	v, err := sv2.NewVersion(s)
	if err != nil {
		util.PrintErr(err, "failed to create semantic version from: %v\n", s)
	}
	c.root.Version = fmt.Sprintf("%s", v)
}

func (c *Commander) Execute() error {
	return c.root.Execute()
}
