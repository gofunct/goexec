package fs

import (
	"bufio"
	"fmt"
	"github.com/gofunct/goexec/pkg/util"
	"github.com/spf13/viper"
	"os"
)

// Prompt prompts user for input with default value.
func (f *Fs) Prompt(key, question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
	text, _ := reader.ReadString('\n')
	f.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (f *Fs) PromptBool(key string, question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("bool | y/n | " + question)
	text, _ := reader.ReadString('\n')
	ans := util.AsBool(text)
	viper.SetDefault(key, ans)
	return ans
}

func enquire(key string) (string, string) {
	return key, fmt.Sprintf("required | please set %s:", key)
}
func (f *Fs) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (f *Fs) Println(msg string) {
	fmt.Println(msg)
}

func (f *Fs) Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (f *Fs) Panic(err error, msg string) {
	if err != nil {
		f.Println(msg)
		panic(err.Error())
	}
}

func (f *Fs) PrintErr(err error, msg string) {
	if err != nil {
		f.Println(err.Error())
		f.Println(msg)
	}
}
