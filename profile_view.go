package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ViewConfig struct {
	filePath    string
	profileName string
	entries     int
	lastMod     time.Time
}

type Viewer struct {
	app      *tview.Application
	config   ViewConfig
	textView *tview.TextView
	header   *tview.TextView
	status   *tview.TextView
}

func NewViewer(config ViewConfig) *Viewer {
	app := tview.NewApplication()
	viewer := &Viewer{
		app:      app,
		config:   config,
		textView: tview.NewTextView(),
		header:   tview.NewTextView(),
		status:   tview.NewTextView(),
	}

	viewer.textView.SetDynamicColors(true)
	viewer.textView.SetRegions(true)
	viewer.textView.SetScrollable(true)
	viewer.textView.SetWrap(false)

	viewer.header.SetTextAlign(tview.AlignLeft)
	viewer.header.SetDynamicColors(true)
	viewer.header.SetBackgroundColor(tcell.ColorDefault)

	viewer.status.SetTextAlign(tview.AlignCenter)
	viewer.status.SetDynamicColors(true)
	viewer.status.SetBackgroundColor(tcell.ColorDefault)

	return viewer
}

func (v *Viewer) layout() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	flex.AddItem(v.header, 2, 0, false)

	flex.AddItem(v.textView, 0, 1, true)

	flex.AddItem(v.status, 1, 0, false)

	return flex
}

func (v *Viewer) updateHeader() {
	headerText := fmt.Sprintf(
		"[yellow]Profile:[white] %s  [yellow]Entries:[white] %d  [yellow]Last Modified:[white] %s",
		v.config.profileName,
		v.config.entries,
		v.config.lastMod.Format("2006-01-02 15:04:05"),
	)
	v.header.SetText(headerText)
}

func (v *Viewer) updateStatus() {
	statusText := "[yellow]↑↓:[white] Scroll  [yellow]PgUp/PgDn:[white] Page Scroll  [yellow]q:[white] Quit"
	v.status.SetText(statusText)
}

func (v *Viewer) highlightContent(content string) string {
	var highlighted strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			highlighted.WriteString("[green]" + line + "[white]\n")
		} else if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				highlighted.WriteString(fmt.Sprintf("[yellow]%s[white]=%s\n", parts[0], parts[1]))
			} else {
				highlighted.WriteString(line + "\n")
			}
		} else {
			highlighted.WriteString(line + "\n")
		}
	}

	return highlighted.String()
}

func (v *Viewer) Run() error {
	content, err := os.ReadFile(v.config.filePath)
	if err != nil {
		return fmt.Errorf("failed to read profile: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	entryCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			entryCount++
		}
	}
	v.config.entries = entryCount

	fileInfo, err := os.Stat(v.config.filePath)
	if err == nil {
		v.config.lastMod = fileInfo.ModTime()
	}

	v.updateHeader()
	v.updateStatus()
	v.textView.SetText(v.highlightContent(string(content)))

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyRune:
			if event.Rune() == 'q' {
				v.app.Stop()
				return nil
			}
		case tcell.KeyPgUp:
			row, _ := v.textView.GetScrollOffset()
			v.textView.ScrollTo(row-10, 0)
			return nil
		case tcell.KeyPgDn:
			row, _ := v.textView.GetScrollOffset()
			v.textView.ScrollTo(row+10, 0)
			return nil
		}
		return event
	})

	if err := v.app.SetRoot(v.layout(), true).Run(); err != nil {
		return fmt.Errorf("failed to start viewer: %v", err)
	}

	return nil
}

func ViewProfile(name string) error {
	profilePath := fmt.Sprintf("/path/to/profiles/%s.env", name)

	config := ViewConfig{
		filePath:    profilePath,
		profileName: name,
	}

	viewer := NewViewer(config)
	return viewer.Run()
}
