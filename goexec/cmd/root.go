// Copyright © 2019 Coleman Word coleman.word@gofunct.com
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

package cmd

import (
	"fmt"
	"github.com/gofunct/fsctl"
	"os"

	"github.com/spf13/cobra"
)

var (
	Fs *fsctl.Fs
)

func init() {
	Fs = fsctl.NewFs()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goexec",
	Short: "goexec is a scripting utility tool that combines the power of golang and bash scripts as a single binary",
	Long: `
Author: Coleman Word coleman.word@gofunct.com
Download: go get github.com/gofunct/goexec/goexec
License: MIT
  _____       ______               
 / ____|     |  ____|              
| |  __  ___ | |__  __  _____  ___ 
| | |_ |/ _ \|  __| \ \/ / _ \/ __|
| |__| | (_) | |____ >  <  __/ (__ 
 \_____|\___/|______/_/\_\___|\___|
                                   

## Getting Started
- Run goexec init to create a default goexec sripting cli in your current directory
- You will be prompted for the name of the directory to create(this will be the name of the cli)
- After entering a directory name, a folder containing the cli, along with a default goexec.yaml will be created
- Enter the created directory and install the program
- Run the command to see the default subcommands and flags

The default generated cli will print its usage if called without any subcommands. It should look like this:

## Default Goexec Usage
A goexec cli comes with default subcommands built in.

Here is the usage for a goexec cli that was generated with the name "scripter":

---> Change me to a longer descripton of your goexec cli <---

Usage:
  scripter [command]

Available Commands:
  debug       debug flags or current configuration
  help        Help about any command
  load        download all of your executable dependencies
  script      run bash scripts

Flags:
  -h, --help      help for scripter
      --version   version for scripter

Use "scripter [command] --help" for more information about a command.

"load" subcommands:

download all of your executable dependencies

Usage:
  scripter load [command]

Available Commands:
  get         intelligently download folders/files from a remote url
  mod         download go mod executable dependencies

Flags:
  -h, --help   help for load

"debug" subcommands:
debug flags or current configuration

Usage:
  scripter debug [command]

Available Commands:
  config      debug configuration settings
  flags       debug current flag settings

Flags:
  -h, --help   help for debug

## Configuration

- Your GoExec cli will read in all files that are named either goexec.yaml or goexec.json
- Your configuration settings are automatically synced to environmental variables
- Your configuration settings can be used to template scripts added to the cli

## Templating

- Goexec uses the excellent sprig funcmap functions package ref:https://github.com/Masterminds/sprig
- Goexec renders scripts that contain "{{" with your current configuration settings and the sprig funcmap
ex: 
if you set name: "Coleman Word" in a goexec.yaml file in your current directory and add the following script:
echo "hello {{ .name }}" >> name.txt

it will create a name.txt file containing "Coleman Word"

## Contributing

Feel free to open up any prs if you find any bugs!

"

`,
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}