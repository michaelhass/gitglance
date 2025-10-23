package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog/label"
	styles "github.com/michaelhass/gitglance/internal/core/ui/style"
)

const (
	titleHeight        = 1
	borderPadding  int = 1
	messagePadding int = 1
	borderWidth    int = 1
)

var (
	titleStyle   = styles.Title.Height(titleHeight)
	borderStyle  = styles.FocusBorder.PaddingLeft(borderPadding).PaddingRight(borderPadding)
	messageStyle = styles.Text.PaddingTop(messagePadding).PaddingBottom(messagePadding)
)

type Model struct {
	title        string
	messageLabel label.MultiLine

	confirmCmd tea.Cmd
	keys       KeyMap

	width, maxContentWidth   int
	height, maxContentHeight int
}

type confirmExecutedMsg struct{}

func New(title string, message string, confirmCmd tea.Cmd) Model {
	return Model{
		title:        title,
		messageLabel: label.NewDefaultMultiLine().SetText(message),
		confirmCmd:   confirmCmd,
		keys:         NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, m.keys.confirm) {
		return m, tea.Sequence(m.confirmCmd, func() tea.Msg { return confirmExecutedMsg{} })
	}
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

	return borderStyle.
		Render(content)
}

func (m Model) Help() []key.Binding {
	return []key.Binding{
		m.keys.confirm,
		m.keys.cancel,
	}
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
