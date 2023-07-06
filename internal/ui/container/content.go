package container

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Content interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	UpdateFocus(isFocused bool) (Content, tea.Cmd)
	View() string
	Title() string
	SetSize(width, height int) Content
	KeyMap() help.KeyMap
}
