package dialog

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/style"
)

type DisplayMode byte

const (
	CenterDisplayMode DisplayMode = iota
	FullScreenDisplayMode
)

const helpHeight = 1

var helpStyle = style.ShortHelp

type Model struct {
	content                         Content
	help                            help.Model
	onCloseCmd                      tea.Cmd
	keys                            key.Binding
	width, height, maxContentHeight int
	displayMode                     DisplayMode
}

func New(content Content, onCloseCmd tea.Cmd, displayMode DisplayMode) Model {
	help := help.NewModel()
	help.ShowAll = false

	return Model{
		content:     content,
		help:        help,
		onCloseCmd:  onCloseCmd,
		displayMode: displayMode,
	}
}

func (m Model) Init() tea.Cmd {
	return m.content.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		return m, tea.Sequence(m.onCloseCmd, Close())
	}

	if _, ok := msg.(CloseMsg); ok {
		return m, m.onCloseCmd
	}

	content, cmd := m.content.Update(msg)
	m.content = content

	return m, cmd
}

func (m Model) View() string {
	content := lipgloss.Place(
		m.width, m.height-helpHeight,
		lipgloss.Center, lipgloss.Center,
		m.content.View(),
		lipgloss.WithWhitespaceBackground(lipgloss.NoColor{}),
	)
	help := helpStyle.Render(m.help.View(m))
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

func (m Model) SetSize(width, height int) Model {
	m.width, m.height = width, height
	m.maxContentHeight = height - helpHeight
	m.help.Width = width - helpStyle.GetHorizontalMargins()

	switch m.displayMode {
	case CenterDisplayMode:
		m.content = m.content.SetSize(width/2, m.maxContentHeight/2) // - margin
	case FullScreenDisplayMode:
		m.content = m.content.SetSize(width, m.maxContentHeight)
	}
	return m
}

func (m Model) OnCloseCmd() tea.Cmd {
	return m.onCloseCmd
}

func (m Model) ShortHelp() []key.Binding {
	return m.content.Help()
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
