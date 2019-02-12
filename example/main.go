package main

import (
	"fmt"
	"github.com/gofunct/goexec"
	"net/http"
	"os"
)

func init() {
	exe.Flags().IntVar(&port, "port", 8080, "port to listen on")
}

var (
	port int
	exe      = goexec.NewCommand("example", "just an example", "0.1")
)

func main() {
	exe.Act("hello", "just sayin hello", func(cmd *goexec.Command) error {
		cmd.AddScript(`echo "hello, {{ .user }} >> hello.txt"`)
		return cmd.Run()
	})
	exe.Act("serve", "serve commands over http", func(cmd *goexec.Command) error {
		cmd.AddScript(`echo "{{ .user }}!"`)
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			cmd.SetOutput(w)
			cmd.Execute()
		})
		return http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	})

	if err := exe.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
