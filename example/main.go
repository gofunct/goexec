package main

import (
	"fmt"
	"github.com/gofunct/goexec"
	"os"
)

func init() {
	exe.Flags().StringVar(&variable, "var", "hello dude", "just an example variable that can be set")
}

var (
	variable string
	exe      = goexec.NewCommand("example", "just an example", "0.1")
)

func main() {
	exe.Act("hello", "just sayin hello", func(cmd *goexec.Command) error {
		cmd.AddScript(`echo "hello, {{ .user }}" >> ./output/example.txt`)
		return cmd.Run()
	})
	exe.Act("dude", "greeting a fellow dude", func(cmd *goexec.Command) error {
		cmd.AddScript(`echo "{{ .var }}!" >> ./output/dude.txt`)
		return cmd.Run()
	})
	if err := exe.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
