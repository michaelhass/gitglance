package textinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/style"
)

var (
	countStyle = style.SublteText.Height(1)
)

type Model struct {
	title string

	textarea textarea.Model

	width  int
	height int

	isFocused bool
}

func New(title string, placeholder string) Model {
	var (
		textarea     = textarea.New()
		blurredStyle = textarea.BlurredStyle
		focusedStyle = textarea.FocusedStyle
	)

	textarea.Placeholder = placeholder
	textarea.Prompt = ""
	textarea.ShowLineNumbers = false

	blurredStyle.Text = style.SublteText
	blurredStyle.Placeholder = style.SublteText
	blurredStyle.CursorLine = lipgloss.NewStyle()
	textarea.BlurredStyle = blurredStyle

	focusedStyle.Placeholder = style.SublteText
	focusedStyle.Text = style.Text
	focusedStyle.CursorLine = lipgloss.NewStyle()
	textarea.FocusedStyle = focusedStyle

	return Model{
		title:    title,
		textarea: textarea,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	textarea, cmd := m.textarea.Update(msg)
	m.textarea = textarea
	return m, cmd
}

func (m Model) UpdateFocus(isFocused bool) (Model, tea.Cmd) {
	m.isFocused = isFocused
	if isFocused {
		m.textarea.Focus()
	} else {
		m.textarea.Blur()
	}
	return m, nil
}

func (m Model) View() string {
	inputLength := len([]rune(m.textarea.Value()))

	var countLine string
	if inputLength > 0
		count := countStyle.
			MaxWidth(m.width - 2).
			Render(fmt.Sprintf("[%d chars]", inputLength))

		countLine = lipgloss.PlaceHorizontal(
			m.width-2,
			lipgloss.Right,
			count,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.textarea.View(),
		countLine,
	)
}

func (m Model) Title() string {
	return m.title
}

func (m Model) SetValue(value string) Model {
	m.textarea.SetValue(value)
	return m
}

func (m Model) SetCursorToStart() Model {
	// textarea.Line() does not seem to return the correct current
	// line of the cursor. At least after setting a new value.
	// Thus, move the cursor up more than potentially needed to ensure
	// that we are at the very beginning of the text input.
	for i := 0; i < m.textarea.LineCount(); i++ {
		m.textarea.CursorUp()
	}
	m.textarea.SetCursor(0) // Only moves to the beginning of the row
	return m
}

func (m Model) SetSize(width, height int) Model {
	m.width, m.height = width, height
	m.textarea.SetWidth(width - 2)
	m.textarea.SetHeight(height - 1)
	return m
}

func (m Model) KeyMap() help.KeyMap {
	return nil
}

func (m Model) Text() string {
	return m.textarea.Value()
}
