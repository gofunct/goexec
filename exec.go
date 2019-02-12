package goexec

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/robfig/cron"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func init() {
	for _, i := range initializers {
		i()
	}
}

func OnInitialize(f ...func()) {
	initializers = append(initializers, f...)
}

var initializers []func()

type ActFunc func(cmd *Command) error
type RunFunc func(name, usg string, act ActFunc)

type Command struct {
	*afero.Afero
	v         *viper.Viper
	exe       *exec.Cmd
	Q         *input.UI
	cfgFile   string
	envPrefix string
	dir       string
	cronz     *cron.Cron
	flags     *cobra.Command
	dkr       *client.Client
}

func NewCommand(name string, usage string, version string) *Command {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	cmd := &Command{
		flags: &cobra.Command{
			Use:     name,
			Short:   usage,
			Version: version,
		},
		Afero: &afero.Afero{
			Fs: afero.NewOsFs(),
		},
		cronz: cron.New(),
		dkr:   cli,
	}

	cmd.Flags().StringVar(&cmd.cfgFile, "config", "goexec.yaml", "path to config file")
	cmd.Flags().StringVar(&cmd.dir, "dir", ".", "directory to execute in")
	cmd.Flags().StringVar(&cmd.envPrefix, "envprefix", "", "prefix to environmental variables")
	cmd.v = viper.New()
	cmd.v.SetFs(cmd.Afero)
	cmd.v.AutomaticEnv()
	debug := &cobra.Command{
		Use:   "debug",
		Short: "debug flags or current configuration",
	}
	debug.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "debug current configuration",
		Run: func(_ *cobra.Command, args []string) {
			cmd.v.Debug()
		},
	}, &cobra.Command{
		Use:   "flags",
		Short: "debug flags",
		Run: func(_ *cobra.Command, args []string) {
			cmd.flags.DebugFlags()
		},
	}, &cobra.Command{
		Use:   "cron",
		Short: "list cron entries",
		Run: func(_ *cobra.Command, args []string) {
			for i, e := range cmd.cronz.Entries() {
				cmd.Printf("%v%s%s", i, e.Prev, e.Next)
			}
		},
	})
	cro := &cobra.Command{
		Use:   "cron",
		Short: "start all cron jobs on their schedules",
	}
	cro.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "start the cron scheduler",
		Run: func(_ *cobra.Command, args []string) {
			cmd.cronz.Run()
		},
	})

	cmd.flags.AddCommand(debug, cro)

	_ = cmd.v.BindPFlags(cmd.Flags())
	if cmd.cfgFile != "" {
		cmd.v.SetConfigFile(cmd.cfgFile)
	}
	if cmd.envPrefix != "" {
		cmd.v.SetEnvPrefix(cmd.envPrefix)
	}
	cmd.Sync()
	c := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"bash", "-c"},
		Env:    os.Environ(),
		Dir:    cmd.dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	cmd.exe = c
	if err := cmd.v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cmd.v.ConfigFileUsed())
	} else {
		cmd.PrintErr(err, "failed to read in config")
	}
	return cmd
}

func (c *Command) Execute() error {
	return c.flags.Execute()
}

func (c *Command) Act(name string, usg string, action ActFunc) {
	c.flags.AddCommand(&cobra.Command{
		Use:   name,
		Short: usg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return action(c)
		},
	})
	_ = c.v.BindPFlags(c.flags.PersistentFlags())
	c.Sync()
}

func (c *Command) Flags() *pflag.FlagSet {
	return c.flags.PersistentFlags()
}

func (c *Command) Run() error {
	return c.exe.Run()
}

func (c *Command) GetReader() io.Reader {
	return c.exe.Stdin
}

func (c *Command) MultiRead(r io.Reader) {
	c.exe.Stdin = io.MultiReader(c.exe.Stdin, r)
}

func (c *Command) AddScript(script string) {
	c.exe.Args = append(c.exe.Args, c.Render(script))
}
func (c *Command) AddScriptFromFile(script string) {
	b, err := ioutil.ReadFile(script)
	if err != nil {
		panic(err)
	}
	script = fmt.Sprintf("%s", b)
	c.exe.Args = append(c.exe.Args, c.Render(script))
}
func (c *Command) AddScriptFromReader(reader io.Reader) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	script := fmt.Sprintf("%s", b)
	c.exe.Args = append(c.exe.Args, c.Render(script))
}

func (c *Command) GetDir() string {
	return c.exe.Dir
}

func (c *Command) SetDir(path string) {
	c.exe.Dir = path
}
