package main

import "fmt"

const (
	Version        = "v0.1.0"
	ProjectName    = "envman"
	configDirName  = ".envman"
	configFileName = "config"
)

const configTemplate = `PROFILE_DIR=/home/%s/.envman/ #Keep the leading slash`

var helpText = fmt.Sprintf(`%s %s
Usage:
  %s [command] [flags]

Available Commands:
  create      Create a new environment profile
  edit        Edit an existing environment profile
  show        Display profile contents
  delete      Delete an environment profile
  list        List all available profiles
  load        Load a profile into current shell

Flags:
  -h, --help    Display help information
  -v, --version Display version information

Use "%s [command] --help" for more information about a command.
`, ProjectName, Version, ProjectName, ProjectName)
