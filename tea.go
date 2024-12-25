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
	w = 64
	h = 64
)

type BoolBoardRow [w]bool
type NumBoardRow  [w]int
type BoolBoard [h]BoolBoardRow
type NumBoard  [h]NumBoardRow

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term       string
	profile    string
	width      int
	height     int
	bg         string
	styles     map[string]lipgloss.Style

	board      BoolBoard
	boardcount NumBoard
	debugmode  bool
}

func IncNeighbors(a NumBoard, x int, y int) NumBoard {
	b := a
	if y > 0 {
		b[y-1][x] += 1
		if x > 0 {
			b[y-1][x-1] += 1
		}
		if x < len(a[0])-1 {
			b[y-1][x+1] += 1
		}
	}
	if y < len(a)-1 {
		b[y+1][x] += 1
		if x > 0 {
			b[y+1][x-1] += 1
		}
		if x < len(a[0])-1 {
			b[y+1][x+1] += 1
		}
	}
	if x > 0 {
		b[y][x-1] += 1
	}
	if x < len(a[0])-1 {
		b[y][x+1] += 1
	}
	return b
}

func InitialModel(pty ssh.Pty, renderer *lipgloss.Renderer, bg string, styles map[string]lipgloss.Style) model {
	b := BoolBoard{}
	br := BoolBoardRow{}
	for i := range br {
		br[i] = false
	}
	for i := range b {
		b[i] = br
	}
	bc := NumBoard{}
	bcr := NumBoardRow{}
	for i := range bcr {
		bcr[i] = 0
	}
	for i := range bc {
		bc[i] = bcr
	}
	m := model{
		term:      pty.Term,
		profile:   renderer.ColorProfile().Name(),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		bg:        bg,
		styles:    styles,

		board:     b,
		boardcount: bc,
		debugmode:  false,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			for y, row := range m.board {
				for x := range row {
					m.boardcount[y][x] = 0
				}
			}
			for y, row := range m.board {
				for x, cell := range row {
					if cell {
						m.boardcount = IncNeighbors(m.boardcount, x, y)
					}
				}
			}
			for y, row := range m.boardcount {
				for x, cell := range row {
					if cell == 3 {
						m.board[y][x] = true
					}
					if cell < 2 || cell > 3 {
						m.board[y][x] = false
					}
				}
			}
			return m, nil
		case "s":
			m.board[4][4] = true
			m.board[5][5] = true
			m.board[5][6] = true
			m.board[6][5] = true
			m.board[6][4] = true
			return m, nil
		case "m":
			m.debugmode = !m.debugmode
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := ""
	if m.debugmode {
		for _, row := range m.boardcount {
			for _, cell := range row {
				s += fmt.Sprintf("%2d", cell)
			}
			s += "\n"
		}
		return m.styles["txt"].Render(s) + "\n\n" + m.styles["quit"].Render("Press 'q' to quit\n")
	}

	for _, row := range m.board {
		for _, cell := range row {
			if cell {
				s += "00"
			} else {
				s += ".."
			}
		}
		s += "\n"
	}
	return m.styles["txt"].Render(s) + "\n\n" + m.styles["quit"].Render("Press 'q' to quit\n")
}
