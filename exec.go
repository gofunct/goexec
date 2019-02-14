package goexec

import "github.com/gofunct/goexec/pkg/commander"

type GoExec struct {
	*commander.Commander
}

func init() {
	for _, i := range initializers {
		i()
	}
}

func OnInitialize(f ...func()) {
	initializers = append(initializers, f...)
}

var initializers []func()

func NewGoExec(name, usg string) *GoExec {
	return &GoExec{
		Commander: commander.NewCommander(name, usg),
	}
}
