package filelist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	up    key.Binding
	down  key.Binding
	enter key.Binding
}

func NewKeyMap(enterHelpText string) KeyMap {
	return KeyMap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("⏎", enterHelpText),
		),
	}
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.up, km.down, km.enter}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
