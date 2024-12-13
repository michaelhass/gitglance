package status

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	up            key.Binding
	down          key.Binding
	left          key.Binding
	right         key.Binding
	commit        key.Binding
	focusUnstaged key.Binding
	focusStaged   key.Binding
	focusDiff     key.Binding

	quit key.Binding

	additionalKeyMap help.KeyMap
}

func newKeyMap() KeyMap {
	return KeyMap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		commit: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("⇧+c", "commit"),
		),
		focusUnstaged: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "To Unstaged"),
		),
		focusStaged: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "To Staged"),
		),
		focusDiff: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "To Diff"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	allKeys := []key.Binding{
		k.focusUnstaged, k.focusStaged, k.focusDiff,
		k.commit,
		k.up, k.down, k.left, k.right,
		k.quit,
	}

	if k.additionalKeyMap == nil {
		return allKeys
	}

	var additionalKeys []key.Binding
	for _, additionalKey := range k.additionalKeyMap.ShortHelp() {
		if containsKey(allKeys, additionalKey) {
			continue
		}
		additionalKeys = append(additionalKeys, additionalKey)
	}

	return append(additionalKeys, allKeys...)
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func containsKey(keyBindings []key.Binding, searchKey key.Binding) bool {
	for _, key := range keyBindings {
		if matchesKey(key, searchKey) {
			return true
		}

	}
	return false
}

// Matches checks if both use the same keys
func matchesKey(lhs key.Binding, rhs key.Binding) bool {
	for _, lhsKey := range lhs.Keys() {
		for _, rhsKey := range rhs.Keys() {
			if lhsKey == rhsKey {
				return true
			}
		}
	}
	return false
}
