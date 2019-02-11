package lg

import (
	"go.uber.org/zap"
)

var l, _ = zap.NewDevelopment()

var sug = l.Sugar()

func init() {
	zap.ReplaceGlobals(l)
}

func FatalIfErr(err error, key string, msg string) {
	if err != nil {
		sug.Fatal(zap.Error(err), zap.String(key, msg))
	}
}

func DebugIfErr(err error, key string, msg string) {
	if err != nil {
		sug.Debug(zap.Error(err), zap.String(key, msg))
	}
}

func WarnIfErr(err error, key string, msg string) {
	if err != nil {
		sug.Warn(zap.Error(err), zap.String(key, msg))
	}
}

func PanicIfErr(err error, key string, msg string) {
	if err != nil {
		sug.Panic(zap.Error(err), zap.String(key, msg))
	}
}

func Debug(key string, msg string) {
	sug.Debug(key, msg)
}

func Sync() error {
	return sug.Sync()
}
