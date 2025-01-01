package main

import "fmt"

const (
	Version        = "v0.1.0"
	ProjectName    = "envman"
	configFileName = "config"
)

const configTemplate = `PROFILE_DIR=/home/%s/.envman/ #Keep the leading slash`

var helpText = fmt.Sprintf(`%s %s
Usage:
  %s [command] [flags]

Available Commands:
  init        Initialize envman in your shell
  profile     Manage environment profiles
  load        Load a profile into the current shell

Profile Subcommands:
  create      Create a new environment profile
  edit        Edit an existing environment profile
  show        Display profile contents
  delete      Delete an environment profile
  list        List all available profiles

Examples:
  # Initialize envman
  $ envman init

  # Profile Management
  $ envman profile create server-test    # Create new profile
  $ envman profile list                  # List all profiles
  $ envman profile show server-test      # Show profile contents
  $ envman profile edit server-test      # Edit existing profile
  $ envman profile delete server-test    # Delete profile

  # Load Profile
  $ envman load server-test             # Load profile into current shell

Flags:
  -h, --help    Display help information
  -v, --version Display version information
`, ProjectName, Version, ProjectName)
