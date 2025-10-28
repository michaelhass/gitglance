package info

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/ui/components/label"
	"github.com/michaelhass/gitglance/internal/core/ui/style"
)

const (
	titleHeight        = 1
	borderPadding  int = 1
	messagePadding int = 1
	borderWidth    int = 1
)

var (
	titleStyle   = style.Title.Height(titleHeight)
	borderStyle  = style.FocusBorder.PaddingLeft(borderPadding).PaddingRight(borderPadding)
	messageStyle = style.Text.PaddingTop(messagePadding).PaddingBottom(messagePadding)
)

type Model struct {
	title        string
	messageLabel label.MultiLine

	width, maxContentWidth   int
	height, maxContentHeight int
}

func New(title, message string) Model {
	return Model{
		title:        title,
		messageLabel: label.NewDefaultMultiLine().SetText(message),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	content := lipgloss.NewStyle().
		MaxWidth(m.maxContentWidth).
		MaxHeight(m.maxContentHeight).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render(m.title),
				messageStyle.Render(m.messageLabel.View()),
			),
		)
	return borderStyle.Render(content)
}

func (m Model) Help() []key.Binding {
	return []key.Binding{}
}

func (m Model) SetSize(width, height int) Model {
	m.width = width - 2
	m.height = height - 2

	borderSize := 2*borderWidth + 2*borderPadding
	m.maxContentWidth = m.width - borderSize
	m.maxContentHeight = m.height - borderSize

	m.messageLabel = m.messageLabel.SetWidth(m.maxContentWidth)

	return m
}
