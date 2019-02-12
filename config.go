package goexec

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

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
	if viper.ConfigFileUsed() != "" {
		name := viper.ConfigFileUsed()
		b, _ := ioutil.ReadFile(name)
		s := fmt.Sprintf("%s", b)
		if strings.Contains(s, " {{") {
			s = c.Render(s)
			r := strings.NewReader(s)
			return c.v.ReadConfig(r)
		}
	}
	c.Println("no config file found")
	return nil
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
