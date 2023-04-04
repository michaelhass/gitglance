package status

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type statusKeyMap struct {
	up    key.Binding
	down  key.Binding
	left  key.Binding
	right key.Binding

	focusUnstaged key.Binding
	focusStaged   key.Binding
	focusDiff     key.Binding

	quit key.Binding

	additionalKeyMap help.KeyMap
}

func newStatusKeyMap() statusKeyMap {
	return statusKeyMap{
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
		focusUnstaged: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "Focus Unstaged"),
		),
		focusStaged: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "Focus Staged"),
		),
		focusDiff: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Focus Diff"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (k statusKeyMap) ShortHelp() []key.Binding {
	allKeys := []key.Binding{
		k.focusUnstaged, k.focusStaged, k.focusDiff,
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

func (k statusKeyMap) FullHelp() [][]key.Binding {
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