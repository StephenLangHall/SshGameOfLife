package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

const (
	host = "localhost"
	port = "6"
	w = 64
	h = 64
	MilisecondFrameWait = 100
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term       string
	profile    string
	width      int
	height     int
	bg         string
	styles     map[string]lipgloss.Style

	board      [h][w]bool
	boardcount [h][w]int
	debugmode  bool
	pause      bool

	curx       int
	cury       int
}

func IncNeighbors(a [h][w]int, x int, y int) [h][w]int {
	b := a
	if y > 0 {
		b[y-1][x] += 1
		if x > 0 {
			b[y-1][x-1] += 1
		} else {
			b[y-1][len(a[0])-1] += 1
		}
		if x < len(a[0])-1 {
			b[y-1][x+1] += 1
		} else {
			b[y-1][0] += 1
		}
	} else {
		b[len(a)-1][x] += 1
		if x > 0 {
			b[len(a)-1][x-1] += 1
		} else {
			b[len(a)-1][len(a[0])-1] += 1
		}
		if x < len(a[0])-1 {
			b[len(a)-1][x+1] += 1
		} else {
			b[len(a)-1][0] += 1
		}
	}
	if y < len(a)-1 {
		b[y+1][x] += 1
		if x > 0 {
			b[y+1][x-1] += 1
		} else {
			b[y+1][len(a[0])-1] += 1
		}
		if x < len(a[0])-1 {
			b[y+1][x+1] += 1
		} else {
			b[y+1][0] += 1
		}
	} else {
		b[0][x] += 1
		if x > 0 {
			b[0][x-1] += 1
		} else {
			b[0][len(a[0])-1] += 1
		}
		if x < len(a[0])-1 {
			b[0][x+1] += 1
		} else {
			b[0][0] += 1
		}
	}
	if x > 0 {
		b[y][x-1] += 1
	} else {
		b[y][len(a[0])-1] += 1
	}
	if x < len(a[0])-1 {
		b[y][x+1] += 1
	} else {
		b[y][0] += 1
	}
	return b
}

func InitialModel(pty ssh.Pty, renderer *lipgloss.Renderer, bg string, styles map[string]lipgloss.Style) model {
	b := [h][w]bool{}
	b = [len(b)][len(b[0])]bool{}
	bc := [h][w]int{}
	bc = [len(bc)][len(bc[0])]int{}
	m := model{
		term:       pty.Term,
		profile:    renderer.ColorProfile().Name(),
		width:      pty.Window.Width,
		height:     pty.Window.Height,
		bg:         bg,
		styles:     styles,

		board:      b,
		boardcount: bc,
		debugmode:  false,
		pause:      true,

		curx:       0,
		cury:       0,
	}
	return m
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(MilisecondFrameWait * time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) UpdateBoard() (tea.Model, tea.Cmd) {
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
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.pause = !m.pause
		case "m":
			m.debugmode = !m.debugmode
		case "w":
			if m.cury > 0 {
				m.cury -= 1
			}
		case "s":
			if m.cury < len(m.board)-1 {
				m.cury += 1
			}
		case "a":
			if m.curx > 0 {
				m.curx -= 1
			}
		case "d":
			if m.curx < len(m.board[0])-1 {
				m.curx += 1
			}
		case "e":
			if m.pause {
				m.board[m.cury][m.curx] = !m.board[m.cury][m.curx]
			}
		case "b":
			SaveBoard(m.board)
		case "c":
			m.boardcount = [len(m.boardcount)][len(m.boardcount[0])]int{}
			m.board = [len(m.boardcount)][len(m.boardcount[0])]bool{}
		case "r":
			m.board = OpenBoard()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case TickMsg:
		if !m.pause {
			m.boardcount = [len(m.boardcount)][len(m.boardcount[0])]int{}
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
		}
		return m, doTick()
	}
	return m, nil
}

func (m model) View() string {
	s := ""
	if m.debugmode {
		for y, row := range m.boardcount {
			for x, cell := range row {
				if y == m.cury && x == m.curx {
					s += m.styles["cur"].Render(fmt.Sprintf("%2d", cell))
				} else {
					s += m.styles["txt"].Render(fmt.Sprintf("%2d", cell))
				}
			}
			s += "\n"
		}
		return m.styles["txt"].Render(s) + "\n\n" + m.styles["quit"].Render("Press 'q' to quit\n")
	}

	for y, row := range m.board {
		for x, cell := range row {
			if cell {
				if y == m.cury && x == m.curx {
					s += m.styles["cur"].Render("<>")
				} else {
					s += m.styles["cur"].Render("  ")
				}
			} else {
				if y == m.cury && x == m.curx {
					s += m.styles["txt"].Render("<>")
				} else {
					s += m.styles["txt"].Render("  ")
				}
			}
		}
		s += "\n"
	}
	return m.styles["txt"].Border(lipgloss.ThickBorder(), true, true, true, true).Render(s) + "\n\n" + m.styles["quit"].Render("Press 'q' to quit\n")
}
