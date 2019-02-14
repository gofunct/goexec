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
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug flags or current configuration",
}

func init() {
	debugCmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "debug current configuration",
		Run: func(_ *cobra.Command, args []string) {
			cmd.v.Debug()
		},
	}, &cobra.Command{
		Use:   "flags",
		Short: "debug flags",
		Run: func(_ *cobra.Command, args []string) {
			cmd.rootcmd.DebugFlags()
		},
	rootCmd.AddCommand(debugCmd)

}
