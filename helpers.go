package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func EnsureConfig() error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	configDir := filepath.Join("/home", currentUser.Username, ".config", ProjectName)
	configFile := filepath.Join(configDir, configFileName)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		content := fmt.Sprintf(configTemplate, currentUser.Username)
		if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
	}

	return nil
}

func (m AppModel) View() string {

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	return helpText
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}
func getEnvmanRoot() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}
	configFilePath := fmt.Sprintf("/home/%s/.config/%s/%s", currentUser.Username, ProjectName, configFileName)
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %v", err)
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
		return "", fmt.Errorf("PROFILE_DIR not found in config")
	}
	return profileDir, nil
}
