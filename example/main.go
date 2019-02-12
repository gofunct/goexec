package main

import (
	"fmt"
	"github.com/gofunct/goexec"
	"os"
	"context"
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
	exe.Act("list-images", "list docker images", func(cmd *goexec.Command) error {

		return cmd.ListImages(context.Background())
	})

	if err := exe.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
