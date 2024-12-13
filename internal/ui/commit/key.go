package commit

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	up          key.Binding
	down        key.Binding
	toggleFocus key.Binding
	commit      key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		toggleFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("⇥", "toggle focus"),
		),
		commit: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("⇧+c", "commit"),
		),
	}
}
