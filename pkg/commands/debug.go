// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gofunct/goexec/pkg/util"
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug flags or current configuration",
}

var cfgDebug = &cobra.Command{
	Use:   "config",
	Short: "debug configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", jsonSettingsString())
	},
}

var flagDebug = &cobra.Command{
	Use:   "flags",
	Short: "debug current flag settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Root().DebugFlags()
	},
}

func init() {
	DebugCmd.AddCommand(cfgDebug, flagDebug)
}

func jsonSettingsString() string {
	return (toPrettyJsonString(util.V.AllSettings()))
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}
