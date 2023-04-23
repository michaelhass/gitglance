package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DisplayMode int

const (
	CenterDisplayMode DisplayMode = iota
	FullScreenDisplayMode
)

type Model struct {
	content     Content
	displayMode DisplayMode
}

func New(content Content, displayMode DisplayMode) Model {
	return Model{
		content:     content,
		displayMode: displayMode,
	}
}

func (m Model) Init() tea.Cmd {
	return m.content.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	content, cmd := m.content.Update(msg)
	m.content = content
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(m.content.View())
}

func (m Model) SetSize(width, height int) Model {
	switch m.displayMode {
	case CenterDisplayMode:
		m.content = m.content.SetSize(width/2, height-10) // - margin
	case FullScreenDisplayMode:
		m.content = m.content.SetSize(width, height)
	}
	return m
}
