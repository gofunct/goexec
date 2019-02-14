package main

import (
	"github.com/gofunct/goexec"
)

func init() {
	exe = goexec.NewGoExec("example", "just an example yo")
	exe.AddVersion("0.1")
	exe.AddScript("fmt", "go format", ".", `go fmt ./...`)
}

var (
	exe *goexec.GoExec
)

func main() {
	defer func() {
		if err := exe.Execute(); err != nil {
			panic(err)
		}
	}()
}
