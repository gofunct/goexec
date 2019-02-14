package main

import (
	"github.com/gofunct/goexec"
)

func init() {
	exe = goexec.NewGoExec("{{ .app }}", "---> Change me to the usage of your goexec cli <---")
	exe.AddDescription(`---> Change me to a longer descripton of your goexec cli <---`)
	exe.AddVersion("0.1")

}

var (
	exe *goexec.GoExec
)

func main() {
	// Add as many scripts as you want here. They will be added as subcommands under {{ .app }}-exec script
	exe.AddScript(
		// The name of the subcommand that will be created
		"fmt",
		// A one line description of the usage of this script
		"recursively go format current directory",
		// The directory the script will run int
		".",
		// The script itself (can be multiline)
		// Scripts can contain templating that will be rendered from your current configuration. Run {{ .app }}-exec debug config to see your current config
		`go fmt ./...`,
	)
	// The cli entrypoint
	if err := exe.Execute(); err != nil {
		panic(err)
	}
}
