package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ANSI
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
	colorReset  = "\033[0m"
)

const (
	iconCheck   = "✓"
	iconX       = "✗"
	iconInfo    = "ℹ"
	iconWarning = "⚠"
)

type CreateProfileModel struct {
	profileName string
	profilePath string
	err         error
	done        bool
}

func (m CreateProfileModel) Init() tea.Cmd {
	return nil
}

func (m CreateProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e":
			editor := "nano"
			if os.Getenv("EDITOR") != "" {
				editor = os.Getenv("EDITOR")
			}
			cmd := exec.Command(editor, m.profilePath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				m.err = fmt.Errorf("failed to open editor: %v", err)
				return m, nil
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m CreateProfileModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s %s%sError:%s %v\n\n",
			iconX,
			colorRed,
			colorBold,
			colorReset,
			m.err,
		)
	}
	if m.done {
		return fmt.Sprintf("\n%s %s%sProfile created:%s %s\n"+
			"%s %s%sUse:%s envman profile edit %s to edit the profile\n"+
			"\nPress 'e' to edit or 'q' to exit\n",
			iconCheck,
			colorGreen,
			colorBold,
			colorReset,
			m.profileName,
			iconInfo,
			colorYellow,
			colorBold,
			colorReset,
			m.profileName,
		)
	}
	return "\n"
}

func CreateProfile(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	name = strings.TrimSpace(name)
	if strings.Contains(name, "/") {
		return fmt.Errorf("profile name cannot contain '/'")
	}

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	configPath := filepath.Join("/home", currentUser.Username, ".config", ProjectName, configFileName)
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var profileDir string
	for _, line := range strings.Split(string(configContent), "\n") {
		if strings.HasPrefix(line, "PROFILE_DIR=") {
			profileDir = strings.TrimPrefix(line, "PROFILE_DIR=")
			profileDir = strings.Split(profileDir, "#")[0]
			profileDir = strings.TrimSpace(profileDir)
			break
		}
	}

	if profileDir == "" {
		return fmt.Errorf("PROFILE_DIR not found in config")
	}

	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %v", err)
	}

	profilePath := filepath.Join(profileDir, name+".env")
	if _, err := os.Stat(profilePath); err == nil {
		return fmt.Errorf("profile '%s' already exists at %s", name, profilePath)
	}

	if err := os.WriteFile(profilePath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create profile file: %v", err)
	}

	model := CreateProfileModel{
		profileName: name,
		profilePath: profilePath,
		done:        true,
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %v", err)
	}

	return nil
}
