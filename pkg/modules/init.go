package modules

import (
	"github.com/gofunct/goexec/pkg/modules/base"
	"github.com/gofunct/goexec/pkg/util"
	"io"
	"strings"
)

func Init() {
	Fs.SetAssetProcessor(base.NewAssetFunc())
	Fs.SetAssetDirProcessor(base.NewAssetDirFunc())
	util.Panic(Fs.MustExecAssets("scripts", createInit()), "failed to execute template on assets\n%s", "")
	createConfig()
}

var defaultConfig = `# Configurations for {{ .app }}-exec must be named goexec.(json or yaml)
# All config files with this name will be loaded from the current working directory regardless of their placement
# Config files may be used to render templated scripts added to {{ .app }}-exec`

func createConfig() {
	defaultConfig = strings.TrimSpace(defaultConfig)
	f, _ := Fs.Create("goexec.yaml")
	_, err := io.WriteString(f, defaultConfig)
	util.Panic(err, "failed to write config to file: %s\n", f.Name())
}