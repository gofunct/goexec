package util

import (
	"fmt"
	"os"
)

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func Println(msg string) {
	fmt.Println(msg)
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func Panic(err error, msg string) {
	if err != nil {
		Println(msg)
		panic(err.Error())
	}
}

func PrintErr(err error, msg string) {
	if err != nil {
		Println(err.Error())
		Println(msg)
	}
}
