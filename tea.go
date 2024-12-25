package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

const (
	host = "localhost"
	port = "6"
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	styles    map[string]lipgloss.Style
}

func InitialModel(pty ssh.Pty, renderer *lipgloss.Renderer, bg string, styles map[string]lipgloss.Style) model {
	return model{
		term:      pty.Term,
		profile:   renderer.ColorProfile().Name(),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		bg:        bg,
		styles:    styles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s", m.term, m.width, m.height, m.bg, m.profile)
	return m.styles["txt"].Render(s) + "\n\n" + m.styles["quit"].Render("Press 'q' to quit\n")
}
