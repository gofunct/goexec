package goexec

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/hashicorp/go-getter"
	"github.com/jessevdk/go-assets"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	*afero.Afero
	v         *viper.Viper
	exe       *exec.Cmd
	Q         *input.UI
	cfgFile   string
	envPrefix string
	dir       string
	cronz     *cron.Cron
	flags     *cobra.Command
}

func NewCommand(name string, usage string, version string) *Command {
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
			cmd.v.AllSettings()
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
				cmd.Println(fmt.Sprintf("%s%s%s", i, e.Prev, e.Next))
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
		Stdout: cmd.flags.OutOrStdout(),
		Stderr: cmd.flags.OutOrStderr(),
	}
	cmd.exe = c
	if err := cmd.v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cmd.v.ConfigFileUsed())
	} else {
		cmd.PrintErr(err, "failed to read in config")
	}
	return cmd
}

type CronSpec int

const (
	HOURLY CronSpec = iota
	DAILY
	WEEKLY
	MONTHLY
	YEARLY
)

func (s CronSpec) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)
	specs := [...]string{
		"HOURLY",
		"DAILY",
		"WEEKLY",
		"MONTHLY",
		"YEARLY"}
	// → `day`: It's one of the
	// values of Weekday constants.
	// If the constant is Sunday,
	// then day is 0.
	// prevent panicking in case of
	// `day` is out of range of Weekday
	if s < HOURLY || s > YEARLY {
		panic(errors.New("unknown cron spec"))
	}
	// return the name of a Weekday
	// constant from the names array
	// above.
	str := specs[s]
	switch str {
	case "HOURLY":
		str = "@hourly "
	case "DAILY":
		str = "@daily"
	case "WEEKLY":
		str = "@weekly"
	case "MONTHLY":
		str = "@monthly"
	case "YEARLY":
		str = "@yearly"
	}
	return str
}

func (c *Command) Cron(spec CronSpec, fn func()) {
	cronz := cron.New()
	c.Panic(cronz.AddFunc(spec.String(), fn), "failed to add cron")
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

func (c *Command) HandlerFunc() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		c.MultiRead(r.Body)
		c.SetOutput(w)
	}
}

func (c *Command) SetOutput(w io.Writer) {
	c.flags.SetOutput(w)
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

func (c *Command) MultiRead(r io.Reader) {
	c.exe.Stdin = io.MultiReader(c.exe.Stdin, r)
}

func (c *Command) OutOrStdOut() io.Writer {
	return c.flags.OutOrStdout()
}

func (c *Command) OutOrStdErr() io.Writer {
	return c.flags.OutOrStderr()
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


func (c *Command) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		c.v.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range c.v.AllSettings() {
		val, ok := v.(string)
		if ok {
			c.PrintErr(os.Setenv(k, val), "failed to bind config to env variable")
		}
	}
}

func (c *Command) ReadInConfig() error {
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
	c.Println(err.Error() + "\nfailed to read comma seperated values from input")
	c.v.SetDefault(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (c *Command) PromptMap(key string, question string) map[string]string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("map | a=b,c=d | " + question)
	text, _ := reader.ReadString('\n')
	txtMap, err := c.AsMap(text)
	c.Println(err.Error() + "\nfailed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
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
		c.Panic(err, "Could not parse template ")
	}

	f, err := c.Create(file.Name())
	if err != nil {
		c.Panic(err, "Could not create file for writing")
	}
	defer f.Close()
	err = tpl.Execute(f, c.v.AllSettings())
	if err != nil {
		c.Panic(err, "Could not execute template")
	}
}

func (c *Command) WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			c.Panic(err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := c.Afero.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, c.v.AllSettings())
		}
		return nil
	}); err != nil {
		c.Panic(err, "failed to walk templates")
	}
}

func (c *Command) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := c.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := c.Afero.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, c.Chmod(dstfile, 0755)
}

func (c *Command) JsonSettings() []byte {
	return (c.toPrettyJson(c.v.AllSettings()))
}

func (c *Command) JsonSettingsString() string {
	return (c.toPrettyJsonString(c.v.AllSettings()))
}

func (c *Command) YamlSettings() []byte {
	bits, err := yaml.Marshal(c.v.AllSettings())
	c.Panic(err, "failed to unmarshal config to yaml")
	return bits
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (c *Command) toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
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
	c.Panic(errors.New(fmt.Sprintf("cannot convert string to bool. valid inputs:\ntrue: %s\nfalse: %s", validBoolT, validBoolF)), "failed to convert string to bool")
	return false
}

func (c *Command) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		c.Panic(err, "failed to render string")
		buf := bytes.NewBuffer(nil)
		c.Panic(t.Execute(buf, c.v.AllSettings()), "failed to render string")
		return buf.String()
	}
	return s
}

func (c *Command) ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}

func (c *Command) Println(msg string) {
	_, err  := fmt.Fprintln(c.OutOrStdErr(), msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Command) Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (c *Command) Panic(err error, msg string) {
	if err != nil {
		c.Println(msg)
		panic(err.Error())
	}
}

func (c *Command) PrintErr(err error, msg string) {
	if err != nil {
		c.Println(err.Error())
		c.Println(msg)
	}
}

func (c *Command) ScanAndReplaceFile(f afero.File, replacements ...string) {
	nm := f.Name()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}
	if err := c.Remove(f.Name()); err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(strings.NewReader(fmt.Sprintf("%s", d)))
	rep := strings.NewReplacer(replacements...)
	var newstr string
	for scanner.Scan() {
		newstr = rep.Replace(scanner.Text())
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	newf, err := c.Create(nm)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.WriteString(newf, newstr)
	c.Panic(err, "failed to write string to new file")
	c.Println("successfully scanned and replaced: " + f.Name())

}

func enquire(key string) (string, string) {
	return key, fmt.Sprintf("required | please set %s:", key)
}

type Mode int

const _Mode_name = "ANYFILEDIR"

const (
	ANY Mode = iota
	FILE
	DIR
)

var _Mode_index = [...]uint8{0, 3, 7, 10}

func (cmd *Command) Load(mode Mode, src, dst string) {
	var moder getter.ClientMode
	switch mode {
	case ANY:
		moder = getter.ClientModeAny
	case FILE:
		moder = getter.ClientModeFile
	case DIR:
		moder = getter.ClientModeDir
	default:
		fmt.Printf("Invalid client mode, must be 'any', 'file', or 'dir': %s", mode.String())
		os.Exit(1)
	}

	// Get the pwd
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting wd: %s", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Build the client
	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  pwd,
		Mode: moder,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		wg.Wait()
		log.Printf("signal %v", sig)
	case <-ctx.Done():
		wg.Wait()
		log.Printf("success!")
	case err := <-errChan:
		wg.Wait()
		log.Fatalf("Error downloading: %s", err)
	}
}

func (i Mode) String() string {
	if i < 0 || i >= Mode(len(_Mode_index)-1) {
		return "Mode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Mode_name[_Mode_index[i]:_Mode_index[i+1]]
}
