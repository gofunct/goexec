package util

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	V = viper.New()
	V.AutomaticEnv()
}

var (
	V *viper.Viper
)

func InitConfig() error {
	if err := filepath.Walk(os.Getenv("PWD"), func(path string, info os.FileInfo, err error) error {
		PrintErr(err, "prevent panic by handling failure accessing a path %q: %v\n", path)
		if info.IsDir() && info.Name() == "vendor" {
			Printf("skipping a dir while finding config files: %+v\n\n", info.Name())
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".json" && info.Name() == "goexec" {
			b, err := ioutil.ReadFile(path)
			Panic(err, "failed to read config file: %v\n", path)
			b, err = json.Marshal(b)
			Panic(err, "failed to unmarshal config file: %v\n", path)
			Panic(V.ReadConfig(bytes.NewBuffer(b)), "failed to read in config: %v\n", path)
		}
		return nil
	}); err != nil {
		return WrapErr(err, "%v", "Failed to register config files ")
	}
	return nil
}
