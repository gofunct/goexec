package gocfg

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/lg"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	fs = &afero.Afero{
		Fs: afero.NewOsFs(),
	}
)

type GoCfg struct {
	Q *input.UI
	v *viper.Viper
	*afero.Afero
}

func New(cfgFile string, envprefix string) *GoCfg {

	g := &GoCfg{v: viper.GetViper(), Q: input.DefaultUI(), Afero: fs}
	g.v.AutomaticEnv()
	if cfgFile != "" {
		g.v.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".temp" (without extension).
		g.v.AddConfigPath(os.Getenv("PWD"))
		g.v.SetConfigName(".config")
		g.v.SetConfigType("yaml")
		lg.WarnIfErr(errors.New("failed to read config file, reading defaults"), "defaults", "path: PWD, name: .config, type: yaml")
	}
	if envprefix != "" {
		g.v.SetEnvPrefix(envprefix)
	}
	g.v.SetFs(g.Afero)
	// If a config file is found, read it in.
	if err := g.v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", g.v.ConfigFileUsed())
	} else {
		lg.WarnIfErr(err, g.v.ConfigFileUsed(), "failed to read in config")
	}
	return g
}

func (g *GoCfg) BindCmd(root *cobra.Command) func() {
	return func() {
		if root.HasAvailableFlags() {
			lg.DebugIfErr(g.v.BindPFlags(root.Flags()), root.Name(), "failed to bind config to Flags()")
		}
		if root.HasAvailablePersistentFlags() {
			lg.DebugIfErr(g.v.BindPFlags(root.PersistentFlags()), root.Name(), "failed to bind config to PersistentFlags()")
		}
		if root.HasAvailableLocalFlags() {
			lg.DebugIfErr(g.v.BindPFlags(root.LocalFlags()), root.Name(), "failed to bind config to LocalFlags()")
			lg.DebugIfErr(g.v.BindPFlags(root.LocalNonPersistentFlags()), root.Name(), "failed to bind config to LocalNonPersistentFlags()")
		}
		if root.HasAvailableInheritedFlags() {
			lg.DebugIfErr(g.v.BindPFlags(root.InheritedFlags()), root.Name(), "failed to bind config to InheritedFlags()")
		}
		if root.HasAvailableSubCommands() {
			for _, cmd := range root.Commands() {

				if cmd.HasAvailableFlags() {
					lg.DebugIfErr(g.v.BindPFlags(cmd.Flags()), cmd.Name(), "failed to bind config to Flags()")
				}
				if cmd.HasAvailablePersistentFlags() {
					lg.DebugIfErr(g.v.BindPFlags(cmd.PersistentFlags()), cmd.Name(), "failed to bind config to PersistentFlags()")
				}
				if cmd.HasAvailableLocalFlags() {
					lg.DebugIfErr(g.v.BindPFlags(cmd.LocalFlags()), cmd.Name(), "failed to bind config to LocalFlags()")
					lg.DebugIfErr(g.v.BindPFlags(cmd.LocalNonPersistentFlags()), cmd.Name(), "failed to bind config to LocalNonPersistentFlags()")
				}
				if cmd.HasAvailableInheritedFlags() {
					lg.DebugIfErr(g.v.BindPFlags(cmd.InheritedFlags()), cmd.Name(), "failed to bind config to InheritedFlags()")
				}

			}
		}
		g.Sync()
	}
}

func (g *GoCfg) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		g.v.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, g := range g.v.AllSettings() {
		val, ok := g.(string)
		if ok {
			lg.DebugIfErr(os.Setenv(strings.ToUpper(k), val), k, "failed to bind "+val)
		}
	}
}

func (g *GoCfg) ReadFrom(reader io.Reader) error {
	return g.v.ReadConfig(reader)
}

func (g *GoCfg) ReadIn() error {
	return g.v.ReadInConfig()
}

func (g *GoCfg) Set(key string, val interface{}) {
	g.Set(key, val)
}

func (g *GoCfg) Get(key string) interface{} {
	return g.v.Get(key)
}
func (g *GoCfg) GetString(key string) string {
	if !g.v.InConfig(key) {
		return g.Prompt(enquire(key))
	}
	return g.v.GetString(key)
}

func (g *GoCfg) GetBool(key string) bool {
	if !g.v.InConfig(key) {
		return g.PromptBool(enquire(key))
	}
	return g.v.GetBool(key)
}

func (g *GoCfg) GetStringSlice(key string, require bool) []string {
	if !g.v.InConfig(key) {
		return g.PromptCSV(enquire(key))
	}
	return g.v.GetStringSlice(key)
}

func (g *GoCfg) GetStringMapString(key string) map[string]string {
	if !g.v.InConfig(key) {
		return g.PromptMap(enquire(key))
	}
	return g.v.GetStringMapString(key)
}

func (g *GoCfg) BindFlagVal(key string, val viper.FlagValue) error {
	return g.v.BindFlagValue(key, val)
}

func (g *GoCfg) BindFPflags(set *pflag.FlagSet) error {
	return g.v.BindPFlags(set)
}

func (g *GoCfg) Unmarshal(obj interface{}) error {
	return g.v.Unmarshal(obj)
}

func (g *GoCfg) EnvPrefix(s string) {
	g.v.SetEnvPrefix(s)
}

func (g *GoCfg) SetDefault(key string, val interface{}) {
	g.v.SetDefault(key, val)
}

// Prompt prompts user for input with default value.
func (g *GoCfg) Prompt(key, question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
	text, _ := reader.ReadString('\n')
	g.v.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptCSV(key string, question string) []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("csv | x, y, z | " + question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := g.AsCSV(text)
	lg.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	g.v.SetDefault(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptMap(key string, question string) map[string]string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("map | a=b,c=d | " + question)
	text, _ := reader.ReadString('\n')
	txtMap, err := g.AsMap(text)
	lg.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	g.v.SetDefault(key, txtMap)
	return txtMap
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptBool(key string, question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("bool | y/n | " + question)
	text, _ := reader.ReadString('\n')
	ans := g.AsBool(text)
	g.v.SetDefault(key, ans)
	return ans
}

// Template reads a go template and writes it to dist given data.
func (g *GoCfg) ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not parse template ")
	}

	f, err := fs.Create(file.Name())
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not create file for writing")
	}
	defer f.Close()
	err = tpl.Execute(f, g.v.AllSettings())
	if err != nil {
		lg.WarnIfErr(err, file.Name(), "Could not execute template")
	}
}

func (g *GoCfg) WalkTemplates(dir string, outDir string) {

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

			f, err := fs.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, g.v.AllSettings())
		}
		return nil
	}); err != nil {
		lg.WarnIfErr(err, dir+" to "+outDir, "failed to walk templates")
	}
}

func (g *GoCfg) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := fs.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := fs.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, fs.Chmod(dstfile, 0755)
}

func (g *GoCfg) JsonSettings() []byte {
	return (g.toPrettyJson(g.v.AllSettings()))
}

func (g *GoCfg) JsonSettingsString() string {
	return (g.toPrettyJsonString(g.v.AllSettings()))
}

func (g *GoCfg) YamlSettings() []byte {
	bits, err := yaml.Marshal(g.v.AllSettings())
	lg.WarnIfErr(err, g.v.ConfigFileUsed(), "failed to unmarshal config to yaml")
	return bits
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (g *GoCfg) toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (g *GoCfg) toPrettyJson(obj interface{}) []byte {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return output
}

func (g *GoCfg) AsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func (g *GoCfg) AsMap(val string) (map[string]string, error) {
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
	for _, g := range arr {
		strings.TrimSpace(g)
		switch {
		case strings.Contains(g, "="):
			kv := strings.Split(g, "=")
			m[kv[0]] = kv[1]
		case strings.Contains(g, ":"):
			kv := strings.Split(g, ":")
			m[kv[0]] = kv[1]
		case strings.Contains(g, ":"):
			kv := strings.Split(g, ":")
			m[kv[0]] = kv[1]
		}
	}
	return m, nil
}

var validBoolT = []string{"Y", "y", "t", "T"}
var validBoolF = []string{"N", "n", "f", "F"}

func (g *GoCfg) AsBool(s string) bool {
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

func (c *GoCfg) Sub(key string) *GoCfg {
	return &GoCfg{
		Q:     c.Q,
		v:     c.v.Sub(key),
		Afero: c.Afero,
	}
}

func (c *GoCfg) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("gocfg").Funcs(sprig.GenericFuncMap()).Parse(s)
		lg.FatalIfErr(err, t.Name(), "failed to render string")
		buf := bytes.NewBuffer(nil)
		lg.FatalIfErr(t.Execute(buf, c.v.AllSettings()), t.Name(), "failed to render string")
		return buf.String()
	}
	return s
}

func (c *GoCfg) ScanAndReplace(r io.Reader, replacements ...string) {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	for scanner.Scan() {
		rep.Replace(scanner.Text())
	}
}

func (c *GoCfg) ScanAndReplaceBytes(r io.Reader, replacements ...string) {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	for scanner.Scan() {
		rep.Replace(fmt.Sprintf("%s", scanner.Bytes()))
	}
}

func enquire(key string) (string, string) {
	return key, fmt.Sprintf("required | please set %s:", key)
}
