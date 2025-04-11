package filelist

import tea "github.com/charmbracelet/bubbletea"

// File handling

type FocusItemMsg struct {
	Item Item
}

type forceFocusUpdateMsg struct{}

func ForceFocusUpdate() tea.Msg {
	return forceFocusUpdateMsg{}
}

type SelectItemMsg struct {
	Item Item
}

type DeleteItemMsg struct {
	Item Item
}

type SelectAllItemMsg struct {
	Items []Item
}

type TopNoMoreItems struct{}

type TopNoMoreFocusableItems struct{}
type BottomNoMoreFocusableItems struct{}

type NoItemsMsg struct{}
