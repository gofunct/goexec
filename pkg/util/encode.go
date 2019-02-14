package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"strings"
)

func JsonSettings() []byte {
	return (toPrettyJson(viper.AllSettings()))
}

func JsonSettingsString() string {
	return (toPrettyJsonString(viper.AllSettings()))
}

func YamlSettings() []byte {
	bits, err := yaml.Marshal(viper.AllSettings())
	Panic(err, "failed to unmarshal config to yaml")
	return bits
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func toPrettyJson(obj interface{}) []byte {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return output
}

func AsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func AsMap(val string) (map[string]string, error) {
	m := make(map[string]string)
	if val == "" {
		return m, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	arr, err := csvReader.Read()
	if err != nil {
		return m, err
	}
	for _, c := range arr {
		strings.TrimSpace(c)
		switch {
		case strings.Contains(c, "="):
			kv := strings.Split(c, "=")
			m[kv[0]] = kv[1]
		case strings.Contains(c, ":"):
			kv := strings.Split(c, ":")
			m[kv[0]] = kv[1]
		case strings.Contains(c, ":"):
			kv := strings.Split(c, ":")
			m[kv[0]] = kv[1]
		}
	}
	return m, nil
}

var validBoolT = []string{"Y", "y", "t", "T"}
var validBoolF = []string{"N", "n", "f", "F"}

func AsBool(s string) bool {
	for _, v := range validBoolT {
		if s == v {
			return true
		}
	}
	for _, v := range validBoolF {
		if s == v {
			return false
		}
	}
	Panic(errors.New(fmt.Sprintf("cannot convert string to bool. valid inputs:\ntrue: %s\nfalse: %s", validBoolT, validBoolF)), "failed to convert string to bool")
	return false
}
