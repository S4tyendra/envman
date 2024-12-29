package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DeleteProfileModel struct {
	profileName string
	confirmed   bool
	err         error
	done        bool
	profilePath string
}

func (m DeleteProfileModel) Init() tea.Cmd {
	return nil
}

func (m DeleteProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if !m.confirmed {
				m.done = true
				m.err = fmt.Errorf("operation cancelled by user")
				return m, tea.Quit
			}
		case "y", "Y":
			if !m.confirmed {
				m.confirmed = true

				if err := os.Remove(m.profilePath); err != nil {
					m.err = fmt.Errorf("failed to delete profile: %v", err)
					return m, tea.Quit
				}

				m.done = true
				return m, tea.Quit
			}
		case "n", "N":
			if !m.confirmed {
				m.done = true
				m.err = fmt.Errorf("operation cancelled by user")
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m DeleteProfileModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s %s%sError:%s %v\n\n",
			iconX,
			colorRed,
			colorBold,
			colorReset,
			m.err,
		)
	}

	if m.done && m.confirmed {
		return fmt.Sprintf("\n%s %s%sProfile deleted:%s %s\n\n",
			iconCheck,
			colorGreen,
			colorBold,
			colorReset,
			m.profileName,
		)
	}

	if m.done && !m.confirmed {
		return fmt.Sprintf("\n%s %s%sOperation cancelled by user%s\n\n",
			iconInfo,
			colorYellow,
			colorBold,
			colorReset,
		)
	}

	if !m.confirmed {
		return fmt.Sprintf("\n%s %s%sDelete profile:%s %s\n"+
			"%s %s%sAre you sure?%s [y/N]: ",
			iconInfo,
			colorYellow,
			colorBold,
			colorReset,
			m.profileName,
			iconInfo,
			colorRed,
			colorBold,
			colorReset,
		)
	}

	return "\n"
}

func DeleteProfile(name string) error {
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

	profilePath := filepath.Join(profileDir, name+".env")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	model := DeleteProfileModel{
		profileName: name,
		profilePath: profilePath,
	}

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run UI: %v", err)
	}

	finalState := finalModel.(DeleteProfileModel)
	if finalState.err != nil {
		return finalState.err
	}

	return nil
}
