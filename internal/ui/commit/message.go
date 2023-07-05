package commit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/style"
)

var (
	countStyle = style.SublteText.Copy().Height(1)
)

type messageContent struct {
	textarea textarea.Model

	width  int
	height int

	isFocused bool
}

func newMessageContent() messageContent {
	var (
		textarea     = textarea.New()
		blurredStyle = textarea.BlurredStyle
		focusedStyle = textarea.FocusedStyle
	)

	textarea.Placeholder = "Commit message"
	textarea.Prompt = ""
	textarea.ShowLineNumbers = false

	blurredStyle.Text = style.SublteText.Copy()
	blurredStyle.Placeholder = style.SublteText.Copy()
	blurredStyle.CursorLine = lipgloss.NewStyle()
	textarea.BlurredStyle = blurredStyle

	focusedStyle.Placeholder = style.SublteText.Copy()
	focusedStyle.Text = style.Text.Copy()
	focusedStyle.CursorLine = lipgloss.NewStyle()
	textarea.FocusedStyle = focusedStyle

	return messageContent{
		textarea: textarea,
	}
}

func (mc messageContent) Init() tea.Cmd {
	return nil
}

func (mc messageContent) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	textarea, cmd := mc.textarea.Update(msg)
	mc.textarea = textarea
	return mc, cmd
}

func (mc messageContent) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	mc.isFocused = isFocused
	if isFocused {
		mc.textarea.Focus()
	} else {
		mc.textarea.Blur()
	}
	return mc, nil
}

func (mc messageContent) View() string {
	inputLength := len([]rune(mc.textarea.Value()))

	count := countStyle.
		MaxWidth(mc.width - 2).
		Render(fmt.Sprint("Chararcters ", inputLength))

	countLine := lipgloss.PlaceHorizontal(
		mc.width-2,
		lipgloss.Right,
		count,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		mc.textarea.View(),
		countLine,
	)
}

func (mc messageContent) Title() string {
	return "Commmit"
}

func (mc messageContent) SetSize(width, height int) container.Content {
	mc.width, mc.height = width, height
	mc.textarea.SetWidth(width - 2)
	mc.textarea.SetHeight(height - 1)
	return mc
}

func (mc messageContent) KeyMap() help.KeyMap {
	return nil
}

func (mc messageContent) message() string {
	return mc.textarea.Value()
}
