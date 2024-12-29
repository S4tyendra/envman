package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EditorConfig struct {
	filePath       string
	profileName    string
	entries        int
	lastMod        time.Time
	sortBy         string
	unsavedChanges bool
}

type Editor struct {
	app        *tview.Application
	config     EditorConfig
	textArea   *tview.TextArea
	header     *tview.TextView
	status     *tview.TextView
	messages   *tview.TextView
	lastBackup string
	hasChanges bool
}

func NewEditor(config EditorConfig) *Editor {
	return &Editor{
		app:        tview.NewApplication(),
		config:     config,
		hasChanges: false,
	}
}

func (e *Editor) layout() *tview.Flex {
	e.header = tview.NewTextView().
		SetDynamicColors(true)

	e.textArea = tview.NewTextArea().
		SetPlaceholder("Enter your environment variables here...")
	e.textArea.SetBorder(true).
		SetTitle(" Editor ").
		SetTitleColor(tcell.ColorGreen).
		SetTitleAlign(tview.AlignLeft)

	e.status = tview.NewTextView().
		SetDynamicColors(true)
	e.messages = tview.NewTextView().
		SetDynamicColors(true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(e.header, 1, 1, false).
		AddItem(e.textArea, 0, 1, true).
		AddItem(e.messages, 1, 1, false).
		AddItem(e.status, 1, 1, false)

	content, err := os.ReadFile(e.config.filePath)
	if err != nil {
		content = []byte("")
	}
	lines := strings.Split(string(content), "\n")

	if e.config.sortBy == "key" {
		lines = sortBy(lines, "key")
	} else if e.config.sortBy == "keylen" {
		lines = sortBy(lines, "keylen")
	} else if e.config.sortBy == "vallen" {
		lines = sortBy(lines, "vallen")
	}
	content = []byte(strings.Join(lines, "\n"))
	e.textArea.SetText(string(content), false)
	e.textArea.SetOffset(0, 0)
	//e.textArea.SetText(string(content), false)
	e.lastBackup = string(content)
	e.updateStatus()

	e.textArea.SetChangedFunc(func() {
		e.hasChanges = true
		e.updateStatus()
	})

	e.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			e.createBackup()
			if err := saveContent(e.config.filePath, e.textArea.GetText()); err != nil {
				e.messages.SetText("[red::b]Error saving: " + err.Error())
			} else {
				e.hasChanges = false
				e.config.unsavedChanges = false
				e.messages.SetText("[green::b]File saved successfully!")
				e.updateStatus()
				go func() {
					time.Sleep(2 * time.Second)
					e.app.QueueUpdateDraw(func() {
						e.messages.SetText("")
					})
				}()
			}
			return nil
		case tcell.KeyCtrlO:
			e.showSortDialog()
			return nil
		//case tcell.KeyCtrlR:
		//	err := ViewProfile(e.config.profileName)
		//	if err != nil {
		//		fmt.Fprintf(os.Stderr, "%s%s%sError:%s %v\n",
		//			colorRed,
		//			colorBold,
		//			iconX,
		//			colorReset,
		//			err,
		//		)
		//		os.Exit(1)
		//	}
		//	return nil
		case tcell.KeyCtrlBackslash:
			e.commentUncommentLine()
			return nil
		case tcell.KeyCtrlX:
			if e.hasChanges {
				modal := tview.NewModal().
					SetText("You have unsaved changes. Save before quitting?").
					AddButtons([]string{"Save", "Don't Save", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						switch buttonIndex {
						case 0:
							e.createBackup()
							saveContent(e.config.filePath, e.textArea.GetText())
							e.app.Stop()
						case 1:
							e.app.Stop()
						default:
							e.app.SetRoot(flex, true)
						}
					})
				e.app.SetRoot(modal, false)
				return nil
			}
			e.app.Stop()
			return nil
		}
		return event
	})

	return flex
}

func (e *Editor) updateStatus() {
	text := e.textArea.GetText()
	entries := 0
	for _, line := range strings.Split(text, "\n") {
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
			entries++
		}
	}
	e.config.unsavedChanges = e.hasChanges
	e.config.entries = entries
	duplicates := e.detectDuplicateKeys()
	duplicatesText := ""
	if len(duplicates) > 0 {
		duplicatesText = "[red::b]Duplicates: " + strings.Join(duplicates, ", ")
	}

	unsavedText := ""
	if e.config.unsavedChanges {
		unsavedText = "[red::b]Unsaved"
	}

	e.header.SetText(fmt.Sprintf("[green::b]Profile:[white] %s | [green]Entries:[white] %d | [green]Last Modified:[white] %s | [green]Path:[white] %s | [green]Sort:[white] %s %s %s",
		e.config.profileName,
		e.config.entries,
		e.config.lastMod.Format("2006-01-02 15:04:05"),
		e.config.filePath,
		e.config.sortBy,
		unsavedText,
		duplicatesText,
	))
	e.status.SetText(fmt.Sprintf("[yellow]Ctrl+S: Save | Ctrl+X: Quit | Ctrl+O: Sort | Ctrl+\\ : Comment/Uncomment"))
}

func (e *Editor) showSortDialog() {
	modal := tview.NewModal().
		SetText("Sort by:").
		AddButtons([]string{"Key", "Key Length", "Value Length", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonIndex {
			case 0:
				e.config.sortBy = "key"
				content, _ := os.ReadFile(e.config.filePath)
				lines := strings.Split(string(content), "\n")
				lines = sortBy(lines, "key")
				e.textArea.SetText(strings.Join(lines, "\n"), true)
			case 1:
				e.config.sortBy = "keylen"
				content, _ := os.ReadFile(e.config.filePath)
				lines := strings.Split(string(content), "\n")
				lines = sortBy(lines, "keylen")
				e.textArea.SetText(strings.Join(lines, "\n"), true)
			case 2:
				e.config.sortBy = "vallen"
				content, _ := os.ReadFile(e.config.filePath)
				lines := strings.Split(string(content), "\n")
				lines = sortBy(lines, "vallen")
				e.textArea.SetText(strings.Join(lines, "\n"), true)
			}
			e.updateStatus()
			e.app.SetRoot(e.layout(), true)
		})
	e.app.SetRoot(modal, false)
}

func (e *Editor) createBackup() error {
	backupPath := e.config.filePath + ".bak"
	content := []byte(e.textArea.GetText())
	err := os.WriteFile(backupPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}
	return nil
}
func highlightLine(text string, cursorY int) string {
	lines := strings.Split(text, "\n")
	for i := range lines {
		if i == cursorY {
			lines[i] = "[yellow]" + lines[i] + "[white]"
		}
	}
	return strings.Join(lines, "\n")
}

func (e *Editor) commentUncommentLine() {
	row, col, _, _ := e.textArea.GetCursor()

	text := e.textArea.GetText()
	lines := strings.Split(text, "\n")

	if row >= len(lines) {
		return
	}

	line := lines[row]
	var modifiedLine string
	var commentPrefix string

	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "#") {
		commentPrefix = "# "
		modifiedLine = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "#"))
	} else {
		commentPrefix = ""
		modifiedLine = "# " + line
	}

	startIndex := 0
	for i := 0; i < row; i++ {
		startIndex += len(lines[i]) + 1
	}
	endIndex := startIndex + len(line)

	e.textArea.Replace(startIndex, endIndex, modifiedLine)

	var newCol int
	if commentPrefix == "" {
		newCol = col + 2
	} else {
		newCol = col - 2
		if newCol < 0 {
			newCol = 0
		}
	}

	e.textArea.Select(startIndex+newCol, startIndex+newCol)

	e.hasChanges = true
	e.updateStatus()
}

func (e *Editor) Run() error {
	e.app.SetRoot(e.layout(), true).EnableMouse(true)
	return e.app.Run()
}

func saveContent(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func sortBy(lines []string, mode string) []string {
	if mode == "key" {
		sort.Slice(lines, func(i, j int) bool {
			keyI := strings.SplitN(lines[i], "=", 2)[0]
			keyJ := strings.SplitN(lines[j], "=", 2)[0]
			return keyI < keyJ
		})
	} else if mode == "keylen" {
		sort.Slice(lines, func(i, j int) bool {
			keyI := strings.SplitN(lines[i], "=", 2)[0]
			keyJ := strings.SplitN(lines[j], "=", 2)[0]
			return len(keyI) < len(keyJ)
		})
	} else if mode == "vallen" {
		sort.Slice(lines, func(i, j int) bool {
			valI := ""
			valJ := ""
			partsI := strings.SplitN(lines[i], "=", 2)
			partsJ := strings.SplitN(lines[j], "=", 2)
			if len(partsI) > 1 {
				valI = partsI[1]
			}
			if len(partsJ) > 1 {
				valJ = partsJ[1]
			}
			return len(valI) < len(valJ)
		})
	}
	return lines
}

func (e *Editor) detectDuplicateKeys() []string {
	lines := strings.Split(e.textArea.GetText(), "\n")
	keyMap := make(map[string][]int)
	var duplicates []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) > 0 {
			key := strings.TrimSpace(parts[0])
			keyMap[key] = append(keyMap[key], i+1)
		}
	}

	for key, lineNumbers := range keyMap {
		if len(lineNumbers) > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s (%s)", key, strings.Join(strings.Split(strings.Trim(fmt.Sprint(lineNumbers), "[]"), " "), ", ")))
		}
	}

	return duplicates
}

func EditProfile(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
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
	fileInfo, err := os.Stat(profilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	content, _ := os.ReadFile(profilePath)
	entries := 0
	for _, line := range strings.Split(string(content), "\n") {
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
			entries++
		}
	}

	config := EditorConfig{
		filePath:       profilePath,
		profileName:    name,
		entries:        entries,
		lastMod:        fileInfo.ModTime(),
		sortBy:         "none",
		unsavedChanges: false,
	}

	editor := NewEditor(config)
	if err := editor.Run(); err != nil {
		return fmt.Errorf("editor error: %v", err)
	}
	return nil
}
