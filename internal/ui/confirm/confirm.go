package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
	styles "github.com/michaelhass/gitglance/internal/ui/style"
)

const (
	titleHeight        = 1
	borderPadding  int = 1
	messagePadding int = 1
)

var (
	titleStyle   = styles.Title.Height(titleHeight)
	borderStyle  = styles.FocusBorder.PaddingLeft(borderPadding).PaddingRight(borderPadding)
	messageStyle = styles.Text.PaddingTop(messagePadding).PaddingBottom(messagePadding)
)

type Model struct {
	title      string
	message    string
	confirmCmd tea.Cmd
	keys       KeyMap

	width  int
	height int
}

func New(title string, message string, confirmCmd tea.Cmd) Model {
	return Model{
		title:      title,
		message:    message,
		confirmCmd: confirmCmd,
		keys:       NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, m.keys.confirm) {
		return m, tea.Sequence(m.confirmCmd, dialog.Close())
	}
	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render(m.title)
	message := messageStyle.Render(m.message)

	return borderStyle.
		MaxWidth(m.width).
		MaxWidth(m.width).
		Render(lipgloss.JoinVertical(lipgloss.Top, title, message))
}

func (m Model) Help() []key.Binding {
	return []key.Binding{
		m.keys.confirm,
		m.keys.cancel,
	}
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
