package goexec

import (
	"github.com/gofunct/gocfg"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
)

type Command struct {
	cmd *exec.Cmd
	dir string
	*gocfg.GoCfg
}

func NewCommand(cfgFile string, reader io.Reader, writer io.Writer) *Command {
	cmd := &Command{
		dir:   os.Getenv("PWD"),
		GoCfg: gocfg.New(cfgFile),
	}
	cmd.Sync()
	c := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"bash", "-c"},
		Env:    os.Environ(),
		Dir:    cmd.GetDir(),
		Stdin:  reader,
		Stdout: writer,
		Stderr: writer,
	}

	cmd.cmd = c

	return cmd

}

func (c *Command) Runnable() bool {
	switch {
	case c.cmd != nil:
		return true
	default:
		return false
	}
}

func (c *Command) GetReader() io.Reader {
	return c.cmd.Stdin
}

func (c *Command) GetStdOut() io.Writer {
	return c.cmd.Stdout
}

func (c *Command) GetStdErr() io.Writer {
	return c.cmd.Stderr
}

func (c *Command) AddScript(script string) {
	c.cmd.Args = append(c.cmd.Args, c.Render(script))
}

func (c *Command) GetDir() string {
	if c.dir == "" {
		c.dir = os.Getenv("PWD")
	}
	return c.dir
}

func (c *Command) SetDir(path string) {
	c.dir = path
}

func (c *Command) Execute() error {
	if c.Runnable() {
		return c.cmd.Run()
	}
	return errors.New("command is not runnable")

}
