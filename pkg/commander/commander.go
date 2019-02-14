package commander

import (
	"fmt"
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
	util.Panic(util.InitConfig(), "failed to initialize configuration: %s", "")
	cmder := &Commander{
		root: &cobra.Command{
			Use:   name,
			Short: usg,
		},
		script: commands.ScriptCmd,
	}
	cmder.root.SetOutput(util.ColoredOut)
	cmder.script.SetOutput(cmder.root.OutOrStderr())
	commands.LoadCmd.SetOutput(cmder.root.OutOrStderr())
	commands.DebugCmd.SetOutput(cmder.root.OutOrStderr())
	cmder.root.AddCommand(cmder.script, commands.DebugCmd, commands.LoadCmd)
	util.V.Set("name", cmder.root.Use)
	util.V.Set("usage", cmder.root.Short)

	for _, c := range cmder.root.Commands() {
		util.V.Set(c.Use+".name", c.Use)
		util.V.Set(c.Use+".usage", c.Short)
		_ = util.V.BindPFlags(c.Flags())
		_ = util.V.BindPFlags(c.PersistentFlags())
	}
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
		Long:  "Script: " + script,
		Run: func(_ *cobra.Command, args []string) {
			util.Panic(cmd.Run(), "failed to run script: %v\n", script)
		},
	})

}

func (c *Commander) AddDescription(s string) {
	c.root.Long = s
	util.V.Set("description", c.root.Long)
}

func (c *Commander) AddVersion(s string) {
	c.root.Version = fmt.Sprintf("%s", s)
	util.V.Set("version", c.root.Version)
}

func (c *Commander) Execute() error {

	return c.root.Execute()
}
