package main

import (
	"github.com/gofunct/goexec"
)

func main() {
	exe := goexec.NewGoExec("example", "just an example yo", "0.1")
	exe.AddScript("fmt", "go format", ".", `go fmt ./...`)
	if err := exe.Execute(); err != nil {
		panic(err)
	}
}
