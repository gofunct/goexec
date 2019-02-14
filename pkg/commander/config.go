package commander

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/goexec/pkg/util"
	"os"
	"strings"
	"text/template"
)

func (c *Commander) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		util.Panic(err, "failed to create template for string: %v", s)
		buf := bytes.NewBuffer(nil)
		util.Panic(t.Execute(buf, util.V.AllSettings()), "failed to render string: %v", s)
		return buf.String()
	}
	return s
}

func (c *Commander) SyncEnv() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		util.V.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range util.V.AllSettings() {
		val, ok := v.(string)
		if ok {
			util.PrintErr(os.Setenv(k, val), "failed to bind config to env variable: %v\n", val)
		}
	}
}
