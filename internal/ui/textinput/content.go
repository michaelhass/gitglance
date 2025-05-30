package textinput

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

type Content struct {
	title string

	textarea textarea.Model

	width  int
	height int

	isFocused bool
}

func NewContent(title string, placeholder string) Content {
	var (
		textarea     = textarea.New()
		blurredStyle = textarea.BlurredStyle
		focusedStyle = textarea.FocusedStyle
	)

	textarea.Placeholder = placeholder
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

	return Content{
		title:    title,
		textarea: textarea,
	}
}

func (c Content) Init() tea.Cmd {
	return nil
}

func (c Content) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	textarea, cmd := c.textarea.Update(msg)
	c.textarea = textarea
	return c, cmd
}

func (c Content) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	c.isFocused = isFocused
	if isFocused {
		c.textarea.Focus()
	} else {
		c.textarea.Blur()
	}
	return c, nil
}

func (c Content) View() string {
	inputLength := len([]rune(c.textarea.Value()))

	count := countStyle.
		MaxWidth(c.width - 2).
		Render(fmt.Sprint("Chararcters ", inputLength))

	countLine := lipgloss.PlaceHorizontal(
		c.width-2,
		lipgloss.Right,
		count,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		c.textarea.View(),
		countLine,
	)
}

func (c Content) Title() string {
	return c.title
}

func (c Content) SetValue(value string) Content {
	c.textarea.SetValue(value)
	return c
}

func (c Content) SetCursorToStart() Content {
	// textarea.Line() does not seem to return the correct current
	// line of the cursor. At least after setting a new value.
	// Thus, move the cursor up more than potentially needed to ensure
	// that we are at the very beginning of the text input.
	for i := 0; i < c.textarea.LineCount(); i++ {
		c.textarea.CursorUp()
	}
	c.textarea.SetCursor(0) // Only moves to the beginning of the row
	return c
}

func (c Content) SetSize(width, height int) container.Content {
	c.width, c.height = width, height
	c.textarea.SetWidth(width - 2)
	c.textarea.SetHeight(height - 1)
	return c
}

func (c Content) KeyMap() help.KeyMap {
	return nil
}

func (c Content) Text() string {
	return c.textarea.Value()
}
