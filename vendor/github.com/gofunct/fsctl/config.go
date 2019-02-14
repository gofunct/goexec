package fsctl

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/fsctl/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (c *Fs) readInConfigFiles() error {
	if err := filepath.Walk(os.Getenv("PWD"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".yaml" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			b, err = yaml.Marshal(b)
			if err != nil {
				panic(err)
			}
			if err := c.ReadConfig(bytes.NewBuffer(b)); err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *Fs) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		util.Panic(err, "failed to render string")
		buf := bytes.NewBuffer(nil)
		util.Panic(t.Execute(buf, c.AllSettings()), "failed to render string")
		return buf.String()
	}
	return s
}

func (c *Fs) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		c.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range c.AllSettings() {
		val, ok := v.(string)
		if ok {
			util.PrintErr(os.Setenv(k, val), "failed to bind config to env variable")
		}
	}
}
