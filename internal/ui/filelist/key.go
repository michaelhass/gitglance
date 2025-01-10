package filelist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	up    key.Binding
	down  key.Binding
	enter key.Binding
	all   key.Binding
}

func NewKeyMap(allText string, enterHelpText string) KeyMap {
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
		all: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", allText),
		),
	}
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.up, km.down, km.all, km.enter}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
