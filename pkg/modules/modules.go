package modules

import (
	"github.com/gofunct/fsctl"
	"github.com/gofunct/goexec/pkg/util"
	"io"
)

func init() {
	Fs = fsctl.NewFs()
}

var (
	Fs *fsctl.Fs
	app string
)

func createInit() io.Writer {
	app = Fs.Prompt("app", util.GreenStringf("please provide a name for your goexec program ---> "))
	err := Fs.MkdirAll(app, 0755)
	util.Panic(err, "failed to create goexec directory: %s\n", app)
	f, err := Fs.Create(app+"/main.go")
	util.PrintErr(err, "failed to create file: %s\n", app+"/main.go")
	return f
}