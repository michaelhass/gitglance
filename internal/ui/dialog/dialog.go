package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DisplayMode byte

const (
	CenterDisplayMode DisplayMode = iota
	FullScreenDisplayMode
)

type Model struct {
	content       Content
	onCloseCmd    tea.Cmd
	width, height int
	displayMode   DisplayMode
}

func New(content Content, onCloseCmd tea.Cmd, displayMode DisplayMode) Model {
	return Model{
		content:     content,
		onCloseCmd:  onCloseCmd,
		displayMode: displayMode,
	}
}

func (m Model) Init() tea.Cmd {
	return m.content.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		return m, Close()
	}

	content, cmd := m.content.Update(msg)
	m.content = content

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.content.View(),
		lipgloss.WithWhitespaceBackground(lipgloss.NoColor{}),
	)
}

func (m Model) SetSize(width, height int) Model {
	m.width, m.height = width, height
	switch m.displayMode {
	case CenterDisplayMode:
		m.content = m.content.SetSize(width/2, height/2) // - margin
	case FullScreenDisplayMode:
		m.content = m.content.SetSize(width, height)
	}
	return m
}

func (m Model) OnCloseCmd() tea.Cmd {
	return m.onCloseCmd
}
