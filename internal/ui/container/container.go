package container

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	styles "github.com/michaelhass/gitglance/internal/ui/Styles"
)

var (
	inactiveTitleStyle = styles.InactiveTitleStyle.Copy().Height(1)
	focusTitleStyle    = styles.TitleStyle.Copy().Height(1)

	inactiveBorderStyle = styles.InactiveBorderStyle.Copy().PaddingLeft(1)
	focusBorderStyle    = styles.FocusBorderStyle.Copy().PaddingLeft(1)
)

type Content interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	View() string
	Title() string
	SetSize(width, height int) Content
	SetIsFocused(isFocused bool) Content
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
	var (
		borderStyle lipgloss.Style
		titleStyle  lipgloss.Style
		title       string
		content     string
	)

	if m.isFocused {
		borderStyle = focusBorderStyle
		titleStyle = focusTitleStyle

	} else {
		borderStyle = inactiveBorderStyle
		titleStyle = inactiveTitleStyle
	}

	title = titleStyle.Render(m.content.Title())
	content = m.content.View()

	return borderStyle.
		Width(m.width).
		Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, title, "", content))
}

func (m Model) SetIsFocused(isFocused bool) Model {
	m.isFocused = isFocused
	m.content = m.content.SetIsFocused(isFocused)
	return m
}

func (m Model) SetSize(width, height int) Model {
	// Substract borders + padding
	m.width, m.height = width-3, height-2
	// Substract title + spacing
	m.content = m.content.SetSize(m.width, m.height-2)
	return m
}

func (m Model) Content() Content {
	return m.content
}

func (m Model) SetContent(content Content) Model {
	m.content = content
	return m
}
