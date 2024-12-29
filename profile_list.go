package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ProfileInfo struct {
	name         string
	size         int64
	lastModified time.Time
	entries      int
}

type ListProfileModel struct {
	profiles []ProfileInfo
	err      error
	done     bool
}

func (m ListProfileModel) Init() tea.Cmd {
	return nil
}

func (m ListProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ListProfileModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s %s%sError:%s %v\n\n",
			iconX,
			colorRed,
			colorBold,
			colorReset,
			m.err,
		)
	}

	if len(m.profiles) == 0 {
		return fmt.Sprintf("\n%s %s%sNo profiles found%s\n"+
			"%s %s%sCreate one with:%s envman profile create <name>\n\n",
			iconInfo,
			colorYellow,
			colorBold,
			colorReset,
			iconInfo,
			colorGreen,
			colorBold,
			colorReset,
		)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n%s %s%sAvailable Profiles:%s (%d total)\n\n",
		iconInfo,
		colorGreen,
		colorBold,
		colorReset,
		len(m.profiles),
	))

	output.WriteString(fmt.Sprintf("%s%s%-20s %-8s %-19s %-19s%s\n",
		colorBold,
		colorYellow,
		"Profile Name",
		"Entries",
		"Last Modified",
		colorReset,
	))

	output.WriteString(fmt.Sprintf("%s%s%s\n",
		colorYellow,
		strings.Repeat("-", 70),
		colorReset,
	))

	for _, p := range m.profiles {
		profileName := strings.TrimSuffix(p.name, ".env")
		output.WriteString(fmt.Sprintf("%s %-20s %s%-8d%s %-19s %-19s\n",
			colorBold,
			profileName,
			colorGreen,
			p.entries,
			colorReset,
			p.lastModified.Format("2006-01-02 15:04"),
		))
	}

	output.WriteString(fmt.Sprintf("\n%s %s%sCommands:%s\n",
		iconInfo,
		colorYellow,
		colorBold,
		colorReset,
	))
	output.WriteString(fmt.Sprintf("  • Use '%senvman profile edit <name>%s' to edit a profile\n", colorBold, colorReset))
	output.WriteString(fmt.Sprintf("  • Use '%senvman profile delete <name>%s' to delete a profile\n", colorBold, colorReset))

	return output.String()
}

func getProfileEntries(filePath string) int {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}

	lines := strings.Split(string(content), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			count++
		}
	}
	return count
}

func ListProfiles() error {
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

	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		model := ListProfileModel{
			profiles: []ProfileInfo{},
		}

		p := tea.NewProgram(model)
		p.Run()
		return nil
	}

	entries, err := os.ReadDir(profileDir)
	if err != nil {
		return fmt.Errorf("failed to read profiles directory: %v", err)
	}

	var profiles []ProfileInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".env") {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			fullPath := filepath.Join(profileDir, entry.Name())
			entryCount := getProfileEntries(fullPath)
			profiles = append(profiles, ProfileInfo{
				name:         entry.Name(),
				size:         info.Size(),
				lastModified: info.ModTime(),
				entries:      entryCount,
			})
		}
	}

	model := ListProfileModel{
		profiles: profiles,
	}

	p := tea.NewProgram(model)
	p.Run()
	return nil
}
