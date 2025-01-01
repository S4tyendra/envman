package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type InitModel struct {
	shell     string
	rcFile    string
	err       error
	selected  string // "a" for append, "c" for copy
	forShell  bool
	succeeded bool
}

func (m InitModel) Init() tea.Cmd {
	return nil
}

func (m InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			if err := m.appendToRC(); err != nil {
				m.err = err
				return m, tea.Quit
			} else {
				m.succeeded = true
				m.selected = "a"
				return m, tea.Quit
			}
		case "c":
			if err := m.copyToClipboard(); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.succeeded = true
			m.selected = "c"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InitModel) View() string {
	if m.forShell {
		switch m.shell {
		case "bash", "zsh":
			return bashInitScript
		case "fish":
			return fishInitScript
		default:
			return fmt.Sprintf("Unsupported shell: %s\n", m.shell)
		}
	}

	if m.err != nil {
		return fmt.Sprintf("%s%s%sError:%s %v\n",
			colorRed,
			colorBold,
			iconX,
			colorReset,
			m.err,
		)
	}

	if m.succeeded {
		if m.selected == "a" {
			return fmt.Sprintf("%s%s%sSuccess:%s Configuration appended to %s\n",
				colorGreen,
				colorBold,
				iconCheck,
				colorReset,
				m.rcFile,
			)
		}
		if m.selected == "c" {
			return fmt.Sprintf("%s%s%sSuccess:%s Configuration copied to clipboard\n",
				colorGreen,
				colorBold,
				iconCheck,
				colorReset,
			)
		}
	}

	currentUser, _ := user.Current()
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%senvman Initialization%s\n\n", colorBold, colorGreen, colorReset))
	sb.WriteString(fmt.Sprintf("Shell: %s\n", m.shell))
	sb.WriteString(fmt.Sprintf("RC File: %s\n", m.rcFile))
	sb.WriteString(fmt.Sprintf("Current User: %s\n", currentUser.Username))
	sb.WriteString(fmt.Sprintf("Current Time (UTC): %s\n\n", currentTime))

	sb.WriteString("The following will be added to your shell configuration:\n\n")
	sb.WriteString(fmt.Sprintf("%s%s%s\n", colorYellow, getEnvmanBlock(m.shell), colorReset))

	sb.WriteString("\nOptions:\n")
	sb.WriteString(fmt.Sprintf("%s[a]%s Append automatically to %s\n", colorBold, colorReset, m.rcFile))
	sb.WriteString(fmt.Sprintf("%s[c]%s Copy to clipboard\n", colorBold, colorReset))
	sb.WriteString(fmt.Sprintf("%s[q]%s Quit\n", colorBold, colorReset))

	return sb.String()
}

func getEnvmanBlock(shell string) string {
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	envmanRoot, err := getEnvmanRoot()
	if err != nil {
		fmt.Println("Error getting envman root directory: ", err)
		os.Exit(1)
	}
	switch shell {
	case "bash", "zsh":
		return fmt.Sprintf(`# START >>>>>>>>>>>> envman Managed [%s] <<<<<<<<<<<<<<<
envman() {
    if [ "$1" = "load" ]; then
        if [ -z "$2" ]; then
            echo "Usage: envman load <profile>" >&2
            return 1
        fi
        if [ ! -f "%s/$2.env" ]; then
             echo "Profile not found: $2" >&2
            return 1
        fi
        source "%s/$2.env"
    else
        command envman "$@"
    fi
}
# END >>>>>>>>>>>> envman Managed <<<<<<<<<<<<<<<`, currentTime, envmanRoot, envmanRoot)
	case "fish":
		return fmt.Sprintf(`# START >>>>>>>>>>>> envman Managed [%s] <<<<<<<<<<<<<<<
function envman
    if [ "$argv[1]" = "load" ]
         if [ -z "$argv[2]" ]
            echo "Usage: envman load <profile>" >&2
            return 1
         end
        if not test -f "%s/$argv[2].env"
             echo "Profile not found: $argv[2]" >&2
            return 1
        end
        source "%s/$argv[2].env"
    else
        command envman $argv
    end
end
# END >>>>>>>>>>>> envman Managed <<<<<<<<<<<<<<<`, currentTime, envmanRoot, envmanRoot)
	default:
		return fmt.Sprintf("Unsupported shell: %s\n", shell)
	}
}

func (m InitModel) appendToRC() error {
	content := getEnvmanBlock(m.shell)
	rcPath := m.rcFile
	if strings.HasPrefix(rcPath, "~/") {
		home := os.Getenv("HOME")
		rcPath = filepath.Join(home, rcPath[2:])
	}

	existing, err := os.ReadFile(rcPath)
	if err == nil && strings.Contains(string(existing), fmt.Sprintf(`envman() {
    if [ "$1" = "load" ]; then`)) {
		return fmt.Errorf("envman configuration already exists in %s", m.rcFile)
	}

	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", m.rcFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + content + "\n"); err != nil {
		return fmt.Errorf("failed to write to %s: %v", m.rcFile, err)
	}
	return nil
}

func (m InitModel) copyToClipboard() error {
	content := getEnvmanBlock(m.shell)
	err := copyToClipboard(content)

	if err != nil {
		return err
	}

	fmt.Printf("\n%s%s%sConfiguration copied to clipboard!%s\n",
		colorGreen,
		colorBold,
		iconCheck,
		colorReset,
	)
	return nil
}

func copyToClipboard(content string) error {
	err := clipboard.WriteAll(content)
	if err != nil {
		if altErr := tryAlternativeClipboard(content); altErr != nil {

			fmt.Printf("\n%s%s%sClipboard access failed!%s\n",
				colorYellow,
				colorBold,
				iconWarning,
				colorReset,
			)
			fmt.Println("\nPlease copy this manually:")
			fmt.Printf("%s%s%s\n",
				colorYellow,
				content,
				colorReset,
			)
			return fmt.Errorf("failed to copy to clipboard: %v", err)
		}
	}
	return nil
}

func tryAlternativeClipboard(content string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		in, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		if _, err := in.Write([]byte(content)); err != nil {
			return err
		}
		if err := in.Close(); err != nil {
			return err
		}
		return cmd.Wait()
	case "linux":

		if xclipCmd := exec.Command("xclip", "-selection", "clipboard"); xclipCmd != nil {
			xclipCmd.Stdin = strings.NewReader(content)
			if err := xclipCmd.Run(); err == nil {
				return nil
			}
		}

		if xselCmd := exec.Command("xsel", "--clipboard", "--input"); xselCmd != nil {
			xselCmd.Stdin = strings.NewReader(content)
			if err := xselCmd.Run(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("no clipboard mechanism available")
	case "windows":
		cmd := exec.Command("clip")
		in, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		if _, err := in.Write([]byte(content)); err != nil {
			return err
		}
		if err := in.Close(); err != nil {
			return err
		}
		return cmd.Wait()
	}
	return fmt.Errorf("unsupported platform")
}

func InitCommand(forShell bool, shell string) error {

	if forShell {

		switch shell {
		case "bash", "zsh":
			fmt.Print(bashInitScript)
		case "fish":
			fmt.Print(fishInitScript)
		default:
			return fmt.Errorf("unsupported shell: %s", shell)
		}
		return nil
	}

	model := InitModel{
		shell:    shell,
		rcFile:   getShellRC(shell),
		forShell: forShell,
	}

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running program: %v", err)
	}

	if m, ok := finalModel.(InitModel); ok && m.err != nil {
		return m.err
	}

	return nil
}

var bashInitScript string
var fishInitScript string

func getShellRC(shell string) string {
	home := os.Getenv("HOME")
	switch shell {
	case "bash":

		if _, err := os.Stat(filepath.Join(home, ".bash_profile")); err == nil {
			return "~/.bash_profile"
		}
		return "~/.bashrc"
	case "zsh":
		return "~/.zshrc"
	case "fish":
		return "~/.config/fish/config.fish"
	default:
		return "~/.profile"
	}
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	switch {
	case shell == "":
		return "bash" // default to bash if we can't detect
	case shell == "/bin/bash" || shell == "/usr/bin/bash":
		return "bash"
	case shell == "/bin/zsh" || shell == "/usr/bin/zsh":
		return "zsh"
	case shell == "/bin/fish" || shell == "/usr/bin/fish":
		return "fish"
	default:
		return "bash" // default to bash for unknown shells
	}
}

func ensureEnvmanDirs() error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	dirs := []string{
		filepath.Join("/home", currentUser.Username, ".envman"),
		filepath.Join("/home", currentUser.Username, ".envman/completions"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

func init() {
	envmanRoot, err := getEnvmanRoot()
	if err != nil {
		fmt.Println("Error getting envman root directory: ", err)
		os.Exit(1)
	}
	bashInitScript = fmt.Sprintf(`envman() {
    if [ "$1" = "load" ]; then
        if [ -z "$2" ]; then
            echo "Usage: envman load <profile>" >&2
            return 1
        fi
         if [ ! -f "%s/$2.env" ]; then
             echo "Profile not found: $2" >&2
            return 1
        fi
        source "%s/$2.env"
    else
        command envman "$@"
    fi
}
`, envmanRoot, envmanRoot)
	fishInitScript = fmt.Sprintf(`function envman
    if [ "$argv[1]" = "load" ]
         if [ -z "$argv[2]" ]
            echo "Usage: envman load <profile>" >&2
            return 1
         end
        if not test -f "%s/$argv[2].env"
             echo "Profile not found: $argv[2]" >&2
            return 1
        end
        source "%s/$argv[2].env"
    else
        command envman $argv
    end
end
`, envmanRoot, envmanRoot)
}
