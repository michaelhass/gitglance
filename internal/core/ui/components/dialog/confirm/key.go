package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	confirm key.Binding
	cancel  key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		confirm: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp("ctrl+y", "confirm"),
		),
		cancel: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String()),
			key.WithHelp("ESC", "cancel"),
		),
	}
}
