package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	All    key.Binding
	Delete key.Binding
}

func NewKeyMap(allText string, enterHelpText string, deleteHelpText string) KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("⏎", enterHelpText),
		),
		All: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", allText),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", deleteHelpText),
		),
	}
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Up, km.Down, km.All, km.Enter, km.Delete}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
