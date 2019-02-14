![gofunct](https://github.com/gofunct/logo/blob/master/white_logo_dark_background.jpg?raw=true)

# GoExec

**Author:** Coleman Word coleman.word@gofunct.com

**Download**: `go get github.com/gofunct/goexec/goexec`

**License:** MIT

```text
 _____       ______               
 / ____|     |  ____|              
| |  __  ___ | |__  __  _____  ___ 
| | |_ |/ _ \|  __| \ \/ / _ \/ __|
| |__| | (_) | |____ >  <  __/ (__ 
 \_____|\___/|______/_/\_\___|\___|
                                   
```

## Getting Started `goexec init`


- Run goexec init to create a default goexec sripting cli in your current directory
- You will be prompted for the name of the directory to create(this will be the name of the cli)
- After entering a directory name, a folder containing the cli, along with a default goexec.yaml will be created
- Enter the created directory and install the program
- Run the command to see the default subcommands and flags

The default generated cli will print its usage if called without any subcommands. It should look like this:
 
## Default Goexec Usage(generated): `"your-goexec-cli"`
A goexec cli comes with default subcommands built in.

Here is the usage for a goexec cli that was generated with the name "scripter":

```text
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

```


### Load Subcommands: `"your-goexec-cli" load`

```text
download all of your executable dependencies

Usage:
  scripter load [command]

Available Commands:
  get         intelligently download folders/files from a remote url
  mod         download go mod executable dependencies

Flags:
  -h, --help   help for load
```


### Debug Subcommands: `"your-goexec-cli" debug`

```text
debug flags or current configuration

Usage:
  scripter debug [command]

Available Commands:
  config      debug configuration settings
  flags       debug current flag settings

Flags:
  -h, --help   help for debug
```


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
