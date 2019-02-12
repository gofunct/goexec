package main

import (
	"fmt"
	"github.com/gofunct/goexec"
	"os"
)

func init() {
	exe.Flags().IntVar(&port, "port", 8080, "port to listen on")
}

var (
	port int
	exe  = goexec.NewCommand("example", "just an example", "0.1")
)

func main() {
	exe.Act("hello", "just sayin hello", func(cmd *goexec.Command) error {
		cmd.AddScript(`echo "hello, {{ .user }}" >> output/hello.txt`)
		return cmd.Run()
	})
	exe.Act("hello", "just sayin hello", func(cmd *goexec.Command) error {
		cmd.list
		return cmd.Run()
	})

	if err := exe.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
