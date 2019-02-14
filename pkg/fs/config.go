package fs

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (f *Fs) ReadInConfigFiles() error {

	if err := filepath.Walk(os.Getenv("PWD"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".yaml" {

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
		c.Panic(err, "failed to render string")
		buf := bytes.NewBuffer(nil)
		c.Panic(t.Execute(buf, c.AllSettings()), "failed to render string")
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
			c.PrintErr(os.Setenv(k, val), "failed to bind config to env variable")
		}
	}
}
