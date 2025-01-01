package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	if err := EnsureConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("%s %s\n", ProjectName, Version)
			return
		}
	}

	if len(os.Args) <= 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		program := tea.NewProgram(AppModel{})
		if _, err := program.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if len(os.Args) > 2 && os.Args[1] == "profile" && os.Args[2] == "create" {
		if len(os.Args) != 4 {
			fmt.Println("Usage: envman profile create <profile-name>")
			return
		}
		if err := CreateProfile(os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%s Error:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 2 && os.Args[1] == "profile" && os.Args[2] == "delete" {
		if len(os.Args) != 4 {
			fmt.Println("Usage: envman profile delete <profile-name>")
			return
		}
		if err := DeleteProfile(os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%s Error:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}

	if len(os.Args) > 2 && os.Args[1] == "profile" && os.Args[2] == "list" {
		if err := ListProfiles(); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%s Error:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 2 && os.Args[1] == "profile" && os.Args[2] == "edit" {
		if len(os.Args) != 4 {
			fmt.Println("Usage: envman profile edit <profile-name>")
			return
		}
		if err := EditProfile(os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%sError:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 2 && os.Args[1] == "profile" && os.Args[2] == "view" {
		var profileName string
		if len(os.Args) == 4 {
			profileName = os.Args[3]
		}
		if err := ViewProfile(profileName); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%sError:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "init" {
		shell := detectShell()
		forShell := false

		if len(os.Args) > 2 && os.Args[2] == "-" {
			forShell = true
			if len(os.Args) > 3 {
				shell = os.Args[3]
			}
		}

		if err := ensureEnvmanDirs(); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%sError:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}

		if err := InitCommand(forShell, shell); err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%sError:%s %v\n",
				colorRed,
				colorBold,
				iconX,
				colorReset,
				err,
			)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Command '%s' not yet implemented or incorrect command\n", os.Args[1])
}
