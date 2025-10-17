package diff

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	up   key.Binding
	down key.Binding
}

func newDiffKeyMap() KeyMap {
	return KeyMap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
