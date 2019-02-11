package gocfg

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/lg"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
	*viper.Viper
	*afero.Afero
}

func New(cfgFile string) *GoCfg {

	g := &GoCfg{Viper: viper.GetViper(), Q: input.DefaultUI(), Afero: fs}
	g.AutomaticEnv()
	if cfgFile != "" {
		g.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".temp" (without extension).
		g.AddConfigPath(os.Getenv("PWD"))
		g.SetConfigName(".config")
		g.SetConfigType("yaml")
	}
	g.SetFs(g.Afero)
	// If a config file is found, read it in.
	if err := g.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", g.ConfigFileUsed())
	} else {
		lg.FatalIfErr(err, g.ConfigFileUsed(), "failed to read in config")
	}
	return g
}

func (g *GoCfg) BindCmd(root *cobra.Command) func() {
	return func() {
		if root.HasAvailableFlags() {
			lg.DebugIfErr(g.BindPFlags(root.Flags()), root.Name(), "failed to bind config to Flags()")
		}
		if root.HasAvailablePersistentFlags() {
			lg.DebugIfErr(g.BindPFlags(root.PersistentFlags()), root.Name(), "failed to bind config to PersistentFlags()")
		}
		if root.HasAvailableLocalFlags() {
			lg.DebugIfErr(g.BindPFlags(root.LocalFlags()), root.Name(), "failed to bind config to LocalFlags()")
			lg.DebugIfErr(g.BindPFlags(root.LocalNonPersistentFlags()), root.Name(), "failed to bind config to LocalNonPersistentFlags()")
		}
		if root.HasAvailableInheritedFlags() {
			lg.DebugIfErr(g.BindPFlags(root.InheritedFlags()), root.Name(), "failed to bind config to InheritedFlags()")
		}
		if root.HasAvailableSubCommands() {
			for _, cmd := range root.Commands() {

				if cmd.HasAvailableFlags() {
					lg.DebugIfErr(g.BindPFlags(cmd.Flags()), cmd.Name(), "failed to bind config to Flags()")
				}
				if cmd.HasAvailablePersistentFlags() {
					lg.DebugIfErr(g.BindPFlags(cmd.PersistentFlags()), cmd.Name(), "failed to bind config to PersistentFlags()")
				}
				if cmd.HasAvailableLocalFlags() {
					lg.DebugIfErr(g.BindPFlags(cmd.LocalFlags()), cmd.Name(), "failed to bind config to LocalFlags()")
					lg.DebugIfErr(g.BindPFlags(cmd.LocalNonPersistentFlags()), cmd.Name(), "failed to bind config to LocalNonPersistentFlags()")
				}
				if cmd.HasAvailableInheritedFlags() {
					lg.DebugIfErr(g.BindPFlags(cmd.InheritedFlags()), cmd.Name(), "failed to bind config to InheritedFlags()")
				}

			}
		}
		g.Sync()
	}
}

func (g *GoCfg) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		g.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, g := range g.AllSettings() {
		val, ok := g.(string)
		if ok {
			lg.DebugIfErr(os.Setenv(strings.ToUpper(k), val), k, "failed to bind "+val)
		}
	}
}

func (g *GoCfg) JsonSettings() []byte {
	return (g.ToPrettyJson(g.AllSettings()))
}

func (g *GoCfg) JsonSettingsString() string {
	return (g.ToPrettyJsonString(g.AllSettings()))
}

func (g *GoCfg) YamlSettings() []byte {
	bits, err := yaml.Marshal(g.AllSettings())
	lg.WarnIfErr(err, g.ConfigFileUsed(), "failed to unmarshal config to yaml")
	return bits
}

// Prompt prompts user for input with default value.
func (g *GoCfg) Prompt(key, question string) string {
	switch {
	case g.InConfig(key):
		return g.GetString(key)
	case os.Getenv(strings.ToUpper(key)) != "":
		return os.Getenv(strings.ToUpper(key))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	g.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptSet(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	_ = os.Setenv(strings.ToUpper(key), text)
	g.Set(key, text)
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptCSV(key string, question string) []string {
	switch {
	case g.InConfig(key):
		return g.GetStringSlice(key)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := g.AsCSV(text)
	lg.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	g.Set(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptSetCSV(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := g.AsCSV(text)
	lg.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	g.Set(key, txtCsv)
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptMap(key string, question string) map[string]string {
	switch {
	case g.InConfig(key):
		return g.GetStringMapString(key)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtMap, err := g.AsMap(text)
	lg.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	return txtMap
}

// Prompt prompts user for input with default value.
func (g *GoCfg) PromptSetMap(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtMap, err := g.AsMap(text)
	lg.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	g.Set(key, txtMap)
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
	err = tpl.Execute(f, g.AllSettings())
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
			return newt.Execute(f, g.AllSettings())
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

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (g *GoCfg) ToPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (g *GoCfg) ToPrettyJson(obj interface{}) []byte {
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

func (c *GoCfg) Render(s string) string {
	t, err := template.New("gocfg").Funcs(sprig.GenericFuncMap()).Parse(s)
	lg.FatalIfErr(err, t.Name(), "failed to render string")
	buf := bytes.NewBuffer(nil)
	lg.FatalIfErr(t.Execute(buf, c.AllSettings()), t.Name(), "failed to render string")
	return buf.String()
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
