package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Content interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	View() string
	SetSize(width, height int) Content
	Help() []key.Binding
}
