package goexec

import (
	"bufio"
	"fmt"
	"os"
)

// Prompt prompts user for input with default value.
func (c *Command) Prompt(key, question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
	text, _ := reader.ReadString('\n')
	c.v.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (c *Command) PromptCSV(key string, question string) []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("csv | x, y, z | " + question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := c.AsCSV(text)
	c.Println(err.Error() + "\nfailed to read comma seperated values from input")
	c.v.SetDefault(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (c *Command) PromptMap(key string, question string) map[string]string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("map | a=b,c=d | " + question)
	text, _ := reader.ReadString('\n')
	txtMap, err := c.AsMap(text)
	c.Println(err.Error() + "\nfailed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	c.v.SetDefault(key, txtMap)
	return txtMap
}

// Prompt prompts user for input with default value.
func (c *Command) PromptBool(key string, question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("bool | y/n | " + question)
	text, _ := reader.ReadString('\n')
	ans := c.AsBool(text)
	c.v.SetDefault(key, ans)
	return ans
}

func enquire(key string) (string, string) {
	return key, fmt.Sprintf("required | please set %s:", key)
}

func (c *Command) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (c *Command) Println(msg string) {
	fmt.Println(msg)
}

func (c *Command) Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (c *Command) Panic(err error, msg string) {
	if err != nil {
		c.Println(msg)
		panic(err.Error())
	}
}

func (c *Command) PrintErr(err error, msg string) {
	if err != nil {
		c.Println(err.Error())
		c.Println(msg)
	}
}
