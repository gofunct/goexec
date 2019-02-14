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
