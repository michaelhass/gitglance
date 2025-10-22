package list

import tea "github.com/charmbracelet/bubbletea"

// File handling

// FocusItemMsg indicates the currently focued item in the list.
type FocusItemMsg struct {
	Item Item
}

type forceFocusUpdateMsg struct{}

// ForceFocusUpdate will force to update and refocus the already focused item.
// It will call the registeredItemHandler of the list.
func ForceFocusUpdate() tea.Msg {
	return forceFocusUpdateMsg{}
}

// SelectItemMsg indicates the selected item in the list. It is produced after
// `enter` key trigger.
type SelectItemMsg struct {
	Item Item
}

// DeleteItemMsg indicates the intent to delete an item in the list.
// The item has not been removed from the list yet. We can return an appropiate tea.Cmd
// from the itemHandler to actually remove the data and update the views.
type DeleteItemMsg struct {
	Item Item
}

// SelectAllItemMsg produced when all items in a list were selected.
type SelectAllItemMsg struct {
	Items []Item
}

// EditItemMsg is an intent to edit an item.
type EditItemMsg struct {
	Item Item
}

type CustomItemMsg struct {
	Item   Item
	KeyMsg tea.KeyMsg
}

type TopNoMoreItems struct{}

type TopNoMoreFocusableItems struct{}
type BottomNoMoreFocusableItems struct{}

type NoItemsMsg struct{}
