package container

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	styles "github.com/michaelhass/gitglance/internal/ui/Styles"
)

var (
	titleStyle = styles.TitleStyle.Copy().Height(1)

	borderStyle      = styles.BorderStyle.Copy()
	focusBorderStyle = styles.FocusBorderStyle.Copy()
)

type Content interface {
	Title() string
	SetSize(width, height int) Content
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	View() string
}

type Model struct {
	content   Content
	width     int
	height    int
	isFocused bool
}

func NewModel(content Content) Model {
	return Model{
		content: content,
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
	title := titleStyle.Render(m.content.Title())

	var builder strings.Builder
	for i := 0; i < m.width; i++ {
		builder.WriteString("w")
	}
	content := lipgloss.NewStyle().Render(builder.String())

	var style lipgloss.Style
	if m.isFocused {
		style = focusBorderStyle
	} else {
		style = borderStyle
	}

	return style.
		Copy().
		Width(m.width).
		Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, title, content))
}

func (m Model) SetIsFocused(isFocused bool) Model {
	m.isFocused = isFocused
	return m
}

func (m Model) SetSize(width, height int) Model {
	// Substract borders
	m.width, m.height = width-2, height-2
	return m
}
