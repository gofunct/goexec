package goexec

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/lg"
	"github.com/jessevdk/go-assets"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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
	fs        *afero.Afero
	v         *viper.Viper
	exe       *exec.Cmd
	Q         *input.UI
	cfgFile   string
	envPrefix string
	dir       string
	reader    io.Reader
	writer    io.Writer

	flags *cobra.Command
}

func NewCommand(name string, usage string, reader io.Reader, writer io.Writer) *Command {

	if reader == nil {
		reader = os.Stdin
	}
	if writer == nil {
		writer = os.Stdout
	}
	cmd := &Command{
		flags: &cobra.Command{
			Use: name,
			Short: usage,
		},
		fs: &afero.Afero{
			Fs: afero.NewOsFs(),
		},
	}
	cmd.flags.SetOutput(writer)
	cmd.Flags().StringVar(&cmd.cfgFile, "config", "", "path to config file")
	cmd.Flags().StringVar(&cmd.dir, "dir", ".", "directory to execute in")
	cmd.Flags().StringVar(&cmd.envPrefix, "envprefix", "", "prefix to environmental variables")
	cmd.v = viper.New()
	cmd.v.SetFs(cmd.fs)
	cmd.v.AutomaticEnv()
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
		Stdin:  reader,
		Stdout: writer,
		Stderr: writer,
	}
	cmd.exe = c
	if err := cmd.v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cmd.v.ConfigFileUsed())
	}
	return cmd
}

func (c *Command) Act(name string, usg string, action ActFunc) {
	c.flags.AddCommand(&cobra.Command{
		Use: name,
		Short: usg,
		RunE: func(cmd *cobra.Command, args []string) error {
		  return action(c)
		},
	})
}

func (c *Command) AddCommand(cmds ...*cobra.Command) {
	c.flags.AddCommand(cmds...)
}
func (c *Command) Flags() *pflag.FlagSet {
	return c.flags.PersistentFlags()
}

func (c *Command) Output() ([]byte, error) {
	return c.exe.Output()
}

func (c *Command) Run() error {
	return c.exe.Run()
}

func (c *Command) GetReader() io.Reader {
	return c.exe.Stdin
}

func (c *Command) GetStdOut() io.Writer {
	return c.exe.Stdout
}

func (c *Command) GetStdErr() io.Writer {
	return c.exe.Stderr
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

func (c *Command) Execute() error {
	_ = c.v.BindPFlags(c.Flags())
	return c.flags.Execute()
}

func (c *Command) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		c.v.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, c := range c.v.AllSettings() {
		val, ok := c.(string)
		if ok {
			lg.DebugIfErr(os.Setenv(strings.ToUpper(k), val), k, "failed to bind "+val)
		}
	}
}

func (c *Command) ReadFrom(reader io.Reader) error {
	return c.v.ReadConfig(reader)
}

func (c *Command) ReadIn() error {
	return c.v.ReadInConfig()
}

func (c *Command) Set(key string, val interface{}) {
	c.Set(key, val)
}

func (c *Command) Get(key string) interface{} {
	return c.v.Get(key)
}
func (c *Command) GetString(key string) string {
	if !c.v.InConfig(key) {
		return c.Prompt(enquire(key))
	}
	return c.v.GetString(key)
}

func (c *Command) GetBool(key string) bool {
	if !c.v.InConfig(key) {
		return c.PromptBool(enquire(key))
	}
	return c.v.GetBool(key)
}

func (c *Command) GetStringSlice(key string, require bool) []string {
	if !c.v.InConfig(key) {
		return c.PromptCSV(enquire(key))
	}
	return c.v.GetStringSlice(key)
}

func (c *Command) GetStringMapString(key string) map[string]string {
	if !c.v.InConfig(key) {
		return c.PromptMap(enquire(key))
	}
	return c.v.GetStringMapString(key)
}

func (c *Command) BindFlagVal(key string, val viper.FlagValue) error {
	return c.v.BindFlagValue(key, val)
}

func (c *Command) BindFPflags(set *pflag.FlagSet) error {
	return c.v.BindPFlags(set)
}

func (c *Command) Unmarshal(obj interface{}) error {
	return c.v.Unmarshal(obj)
}

func (c *Command) EnvPrefix(s string) {
	c.v.SetEnvPrefix(s)
}

func (c *Command) SetDefault(key string, val interface{}) {
	c.v.SetDefault(key, val)
}

// Prompt prompts user for input with default value.
func (c *Command) Prompt(key, question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
	text, _ := reader.ReadString('\n')
	c.v.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (c *Command) PromptCSV(key string, question string) []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("csv | x, y, z | " + question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := c.AsCSV(text)
	lg.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	c.v.SetDefault(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (c *Command) PromptMap(key string, question string) map[string]string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("map | a=b,c=d | " + question)
	text, _ := reader.ReadString('\n')
	txtMap, err := c.AsMap(text)
	lg.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	c.v.SetDefault(key, txtMap)
	return txtMap
}

// Prompt prompts user for input with default value.
func (c *Command) PromptBool(key string, question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("bool | y/n | " + question)
	text, _ := reader.ReadString('\n')
	ans := c.AsBool(text)
	c.v.SetDefault(key, ans)
	return ans
}

// Template reads a go template and writes it to dist given data.
func (c *Command) ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not parse template ")
	}

	f, err := c.fs.Create(file.Name())
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not create file for writing")
	}
	defer f.Close()
	err = tpl.Execute(f, c.v.AllSettings())
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not execute template")
	}
}

func (c *Command) WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			lg.DebugIfErr(err, path, "walkfunc copy error")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := c.fs.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, c.v.AllSettings())
		}
		return nil
	}); err != nil {
		lg.WarnIfErr(err, dir+" to "+outDir, "failed to walk templates")
	}
}

func (c *Command) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := c.fs.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := c.fs.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, c.fs.Chmod(dstfile, 0755)
}

func (c *Command) JsonSettings() []byte {
	return (c.toPrettyJson(c.v.AllSettings()))
}

func (c *Command) JsonSettingsString() string {
	return (c.toPrettyJsonString(c.v.AllSettings()))
}

func (c *Command) YamlSettings() []byte {
	bits, err := yaml.Marshal(c.v.AllSettings())
	lg.WarnIfErr(err, c.v.ConfigFileUsed(), "failed to unmarshal config to yaml")
	return bits
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (c *Command) toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (c *Command) toPrettyJson(obj interface{}) []byte {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return output
}

func (c *Command) AsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func (c *Command) AsMap(val string) (map[string]string, error) {
	m := make(map[string]string)
	if val == "" {
		return m, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	arr, err := csvReader.Read()
	if err != nil {
		return m, err
	}
	for _, c := range arr {
		strings.TrimSpace(c)
		switch {
		case strings.Contains(c, "="):
			kv := strings.Split(c, "=")
			m[kv[0]] = kv[1]
		case strings.Contains(c, ":"):
			kv := strings.Split(c, ":")
			m[kv[0]] = kv[1]
		case strings.Contains(c, ":"):
			kv := strings.Split(c, ":")
			m[kv[0]] = kv[1]
		}
	}
	return m, nil
}

var validBoolT = []string{"Y", "y", "t", "T"}
var validBoolF = []string{"N", "n", "f", "F"}

func (c *Command) AsBool(s string) bool {
	for _, v := range validBoolT {
		if s == v {
			return true
		}
	}
	for _, v := range validBoolF {
		if s == v {
			return false
		}
	}
	panic(errors.New(fmt.Sprintf("cannot convert string to bool. valid inputs:\ntrue: %s\nfalse: %s", validBoolT, validBoolF)))
}

func (c *Command) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("gocfg").Funcs(sprig.GenericFuncMap()).Parse(s)
		lg.FatalIfErr(err, t.Name(), "failed to render string")
		buf := bytes.NewBuffer(nil)
		lg.FatalIfErr(t.Execute(buf, c.v.AllSettings()), t.Name(), "failed to render string")
		return buf.String()
	}
	return s
}

func (c *Command) ScanAndReplace(r io.Reader, replacements ...string) {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	for scanner.Scan() {
		rep.Replace(scanner.Text())
	}
}

func (c *Command) ScanAndReplaceBytes(r io.Reader, replacements ...string) {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	for scanner.Scan() {
		rep.Replace(fmt.Sprintf("%s", scanner.Bytes()))
	}
}

func enquire(key string) (string, string) {
	return key, fmt.Sprintf("required | please set %s:", key)
}
