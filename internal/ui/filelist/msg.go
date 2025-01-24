package filelist

// File handling

type FocusItemMsg struct {
	Item Item
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
