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

## Goexec Usage: `goexec`

<details><summary>show</summary>
<p>

```text

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

p       he...

{{abbrevboth 5 10 "1234 5678 9123"}}:                                               ...5678...

{{initials "First Try"}}:                                                           FT

{{randNumeric 3}}:                                                                  528

{{wrapWith 5 "\t" "Hello World"}}:                                                  Hello	World

{{contains "cat" "catch"}}:                                                         true

{{hasPrefix "cat" "catch"}}:                                                        true

{{cat "hello" "beautiful" "world"}}:                                                hello beautiful world

{{"I Am Henry VIII" | replace " " "-"}}:                                            I-Am-Henry-VIII

{{snakecase "FirstName"}}:                                                          first_name

{{camelcase "http_server"}}:                                                        HttpServer

{{shuffle "hello"}}:                                                                holle

RegExp:

{{regexMatch "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}" "test@acme.com"}}:   true

{{- /*{{regexFindAll "[2,4,6,8]" "123456789"}}*/}}:
{{regexFind "[a-zA-Z][1-9]" "abcd1234"}}:                                           d1

{{regexReplaceAll "a(x*)b" "-ab-axxb-" "${1}W"}}:                                   -W-xxW-

{{regexReplaceAllLiteral "a(x*)b" "-ab-axxb-" "${1}"}}:                             -${1}-${1}-

{{regexSplit "z+" "pizza" -1}}:                                                     [pi a]

{{if contains "cat" "catch"}}yes{{else}}no{{end}}:   yes

{{1 | plural "one anchovy" "many anchovies"}}:       one anchovy

{{2 | plural "one anchovy" "many anchovies"}}:       many anchovies

{{3 | plural "one anchovy" "many anchovies"}}:       many anchovies


## Contributing

Feel free to open up any prs if you find any bugs!

"

Usage:
  goexec [command]

Available Commands:
  help        Help about any command
  init        initialize a new goexec program in your current working directory

Flags:
  -h, --help   help for goexec

Use "goexec [command] --help" for more information about a command.

```

</p>
</details>
 
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
Here is a short list of template function examples that can be embedded in a script:
(Generating uuid, self-signed certificates, rsa private keys, and more can be rendered at runtime but not listed below) 

Semantic Versioning:
`{{ semverCompare "1.2.3" "1.2.3" }}`:  `true`

		`{{ semverCompare "^1.2.0" "1.2.3" }}`: `true`

		`{{ semverCompare "^1.2.0" "2.2.3" }}`: `false`

Deriving Passwords: 
	`{{derivePassword 1 "long" "password" "user" "example.com"}}`:    "ZedaFaxcZaso9*"

		`{{derivePassword 2 "long" "password" "user" "example.com"}}`:    "Fovi2@JifpTupx"

		`{{derivePassword 1 "maximum" "password" "user" "example.com"}}`: "pf4zS1LjCg&LjhsZ7T2~"

		`{{derivePassword 1 "medium" "password" "user" "example.com"}}`:  "ZedJuz8$"

		`{{derivePassword 1 "basic" "password" "user" "example.com"}}`:   "pIS54PLs"

		`{{derivePassword 1 "short" "password" "user" "example.com"}}`:   "Zed5"

		`{{derivePassword 1 "pin" "password" "user" "example.com"}}`:     "6685",

String Manipulation:
{{trim "   hello    "}}:                                                            hello

{{trimAll "$" "$5.00"}}:                                                          5.00

{{trimSuffix "-" "hello-"}}:                                                        hello

{{upper "hello"}}:                                                                  HELLO

{{lower "HELLO"}}:                                                                  hello

{{title "hello world"}}:                                                            Hello World

{{untitle "Hello World"}}:                                                          hello world

{{repeat 3 "hello"}}:                                                               hellohellohello

{{substr 0 5 "hello world"}}:                                                       hello

{{nospace "hello w o r l d"}}:                                                      helloworld

{{trunc 5 "hello world"}}:                                                          hello

{{abbrev 5 "hello world"}}:                                                         he...

{{abbrevboth 5 10 "1234 5678 9123"}}:                                               ...5678...

{{initials "First Try"}}:                                                           FT

{{randNumeric 3}}:                                                                  528

{{wrapWith 5 "\t" "Hello World"}}:                                                  Hello	World

{{contains "cat" "catch"}}:                                                         true

{{hasPrefix "cat" "catch"}}:                                                        true

{{cat "hello" "beautiful" "world"}}:                                                hello beautiful world

{{"I Am Henry VIII" | replace " " "-"}}:                                            I-Am-Henry-VIII

{{snakecase "FirstName"}}:                                                          first_name

{{camelcase "http_server"}}:                                                        HttpServer

{{shuffle "hello"}}:                                                                holle

RegExp:

{{regexMatch "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}" "test@acme.com"}}:   true

{{- /*{{regexFindAll "[2,4,6,8]" "123456789"}}*/}}:
{{regexFind "[a-zA-Z][1-9]" "abcd1234"}}:                                           d1

{{regexReplaceAll "a(x*)b" "-ab-axxb-" "${1}W"}}:                                   -W-xxW-

{{regexReplaceAllLiteral "a(x*)b" "-ab-axxb-" "${1}"}}:                             -${1}-${1}-

{{regexSplit "z+" "pizza" -1}}:                                                     [pi a]

{{if contains "cat" "catch"}}yes{{else}}no{{end}}:   yes

{{1 | plural "one anchovy" "many anchovies"}}:       one anchovy

{{2 | plural "one anchovy" "many anchovies"}}:       many anchovies

{{3 | plural "one anchovy" "many anchovies"}}:       many anchovies



## Contributing

Feel free to open up any prs if you find any bugs!
