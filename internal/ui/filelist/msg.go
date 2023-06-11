package filelist

// File handling

type FocusItemMsg struct {
	Item Item
}

type SelectItemMsg struct {
	Item Item
}

type TopNoMoreItems struct{}

type TopNoMoreFocusableItems struct{}
type BottomNoMoreFocusableItems struct{}
