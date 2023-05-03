package commit

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type messageContent struct {
	width  int
	height int

	isFocused bool
}

func newMessageContent() messageContent {
	return messageContent{}
}

func (mc messageContent) Init() tea.Cmd {
	return nil
}

func (mc messageContent) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	return mc, nil
}

func (mc messageContent) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	mc.isFocused = isFocused
	return mc, nil
}

func (mc messageContent) View() string {
	return lipgloss.
		NewStyle().
		Width(mc.width).
		Height(mc.height).
		Render("Commit message")
}

func (mc messageContent) Title() string {
	return "Commmit"
}

func (mc messageContent) SetSize(width, height int) container.Content {
	mc.width, mc.height = width, height
	return mc
}

func (mc messageContent) KeyMap() help.KeyMap {
	return nil
}
