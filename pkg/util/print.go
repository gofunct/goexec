package util

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
)

var (
	errorStrf  = color.New(color.FgRed).SprintfFunc()
	debugStrf  = color.New(color.FgCyan).SprintfFunc()
	errorf     = color.New(color.FgRed).PrintfFunc()
	printerror = color.New(color.FgRed).PrintFunc()
	debugf     = color.New(color.FgCyan).PrintfFunc()
	BlueString = color.New(color.FgBlue).SprintfFunc()
	ColoredOut = color.Output
	ColoredErr = color.Error
	GreenStringf = color.New(color.FgGreen).SprintfFunc()
)

func Printf(format string, args ...interface{}) {
	debugf(format, args)
}

func Exit(format string, msg ...interface{}) {
	errorf(format, msg)
	os.Exit(1)
}

func Panic(err error, format string, msgs ...interface{}) {
	if err != nil {
		panic(errorStrf(errors.Wrapf(err, format, msgs).Error()))
	}
}

func PrintErr(err error, format string, msgs ...interface{}) {
	if err != nil {
		errorf("ERROR: %s\n", errors.Wrapf(err, format, msgs))
	}
}

func WrapErr(err error, format string, msgs ...interface{}) error {
	return errors.Wrapf(err, format, msgs)
}
