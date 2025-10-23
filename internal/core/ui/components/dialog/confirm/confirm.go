package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/ui/components/label"
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
	textInput    textinput.Model

	onConfirmCmd       tea.Cmd
	onTextInputConfirm func(string) tea.Cmd

	keys KeyMap

	width, maxContentWidth   int
	height, maxContentHeight int
}

type confirmExecutedMsg struct{}

func New(title string, message string) Model {

	return Model{
		title:        title,
		messageLabel: label.NewDefaultMultiLine().SetText(message),
		keys:         NewKeyMap(),
	}
}

func (m Model) WithOnConfirmCmd(cmd tea.Cmd) Model {
	m.onConfirmCmd = cmd
	// reset textinput
	m.textInput = textinput.Model{}
	m.onTextInputConfirm = nil

	return m
}

func (m Model) WithTextInput(placeholder string, onConfirm func(string) tea.Cmd) Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Focus()
	m.onTextInputConfirm = onConfirm
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, m.keys.confirm) {
		var confirmCmd tea.Cmd
		if m.HasTextInput() {
			m.onTextInputConfirm(m.textInput.Value())
		} else {
			confirmCmd = m.onConfirmCmd
		}
		return m, tea.Sequence(confirmCmd, func() tea.Msg { return confirmExecutedMsg{} })
	}
	input, cmd := m.textInput.Update(msg)
	m.textInput = input
	return m, cmd
}

func (m Model) View() string {
	elements := []string{
		titleStyle.Render(m.title),
		messageStyle.Render(m.messageLabel.View()),
	}
	if m.HasTextInput() {
		elements = append(elements, m.textInput.View())
	}

	content := lipgloss.NewStyle().
		MaxWidth(m.maxContentWidth).
		MaxHeight(m.maxContentHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left, elements...))

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
	if m.HasTextInput() {
		m.textInput.Width = m.maxContentWidth
	}

	return m
}

func (m Model) HasTextInput() bool {
	return m.onTextInputConfirm != nil
}
