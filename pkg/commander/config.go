package commander

import (
	"bytes"
	"encoding/json"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/goexec/pkg/util"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (c *Commander) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		util.Panic(err, "failed to create template for string: %v", s)
		buf := bytes.NewBuffer(nil)
		util.Panic(t.Execute(buf, viper.AllSettings()), "failed to render string: %v", s)
		return buf.String()
	}
	return s
}

func (c *Commander) SyncEnv() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		viper.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range viper.AllSettings() {
		val, ok := v.(string)
		if ok {
			util.PrintErr(os.Setenv(k, val), "failed to bind config to env variable: %v\n", val)
		}
	}
}

func (c *Commander) InitConfig() error {
	c.SyncEnv()
	if err := filepath.Walk(os.Getenv("PWD"), func(path string, info os.FileInfo, err error) error {
		util.PrintErr(err, "prevent panic by handling failure accessing a path %q: %v\n", path)
		if info.IsDir() && info.Name() == "vendor" {
			util.Printf("skipping a dir without errors: %+v\n", info.Name())
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".json" && info.Name() == "goexec" {
			b, err := ioutil.ReadFile(path)
			util.Panic(err, "failed to read config file: %v\n", path)
			b, err = json.Marshal(b)
			util.Panic(err, "failed to unmarshal config file: %v\n", path)
			util.Panic(viper.ReadConfig(bytes.NewBuffer(b)), "failed to read in config: %v\n", path)
		}
		return nil
	}); err != nil {
		return util.WrapErr(err, "%v", "Failed to register config files ")
	}
	c.SyncEnv()
	return nil
}
