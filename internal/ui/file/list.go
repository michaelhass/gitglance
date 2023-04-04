package file

import (
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/styles"
)

var (
	itemStyle         = styles.TextSyle.Copy()
	focusedItemStyle  = styles.FocusTextStyle.Copy()
	inactiveItemStyle = styles.InactiveTextStyle.Copy()
)

type ListItemHandler func(msg tea.Msg) tea.Cmd

type List struct {
	items        []ListItem
	visibleItems []ListItem
	itemHandler  ListItemHandler
	keys         KeyMap
	title        string
	width        int
	height       int
	cursor       int
	pageStartIdx int
	isFocused    bool
}

func NewList(title string, itemHandler ListItemHandler, keys KeyMap) List {
	return List{title: title, itemHandler: itemHandler, keys: keys}
}
func (l List) Init() tea.Cmd {
	return nil
}

func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
	if !l.isFocused || len(l.items) == 0 {
		return l, nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keys.up):
			if l.cursor == 0 {
				l.pageStartIdx = l.nextPageStartIdx(-1)
				l.visibleItems = l.updateVisibleItems()
				break
			}
			l.cursor -= 1
			cmd := l.itemHandler(FocusItemMsg{Item: l.visibleItems[l.cursor]})
			cmds = append(cmds, cmd)
		case key.Matches(msg, l.keys.down):
			if l.cursor >= len(l.visibleItems)-1 {
				l.pageStartIdx = l.nextPageStartIdx(1)
				l.visibleItems = l.updateVisibleItems()
				break
			}
			l.cursor += 1
			cmd := l.itemHandler(FocusItemMsg{Item: l.visibleItems[l.cursor]})
			cmds = append(cmds, cmd)
		case key.Matches(msg, l.keys.enter):
			item := l.visibleItems[l.cursor]
			cmd := l.itemHandler(SelectItemMsg{Item: item})
			cmds = append(cmds, cmd)
		}
	default:
		// Check if the curser is out of bounds due to content change.
		if len(l.visibleItems) == 0 || l.cursor < len(l.visibleItems) {
			break
		}
		l.cursor = len(l.visibleItems) - 1
	}

	return l, tea.Batch(cmds...)
}

func (l List) UpdateFocus(isFocused bool) (List, tea.Cmd) {
	var cmd tea.Cmd
	if isFocused && !l.isFocused && len(l.items) > 0 {
		cmd = l.itemHandler(FocusItemMsg{Item: l.visibleItems[l.cursor]})
	}
	l.isFocused = isFocused
	return l, cmd
}

func (l List) View() string {
	var renderedItems = make([]string, len(l.visibleItems))

	for i, item := range l.visibleItems {
		style := itemStyle
		if !l.isFocused {
			style = inactiveItemStyle
		} else if i == l.cursor {
			style = focusedItemStyle
		}

		runes := []rune(item.String())
		runes = runes[:min(len(runes), l.width-1)]
		itemString := string(runes)
		renderedItems[i] = style.Width(l.width).Render(itemString)
	}

	return lipgloss.
		NewStyle().
		Height(l.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedItems...))
}

func (l List) Title() string {
	return l.title
}

func (l List) SetSize(width, height int) List {
	l.width = width
	l.height = height
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l List) KeyMap() help.KeyMap {
	return l.keys
}

func (l List) SetListItems(items []ListItem) List {
	l.items = items
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l List) FocusedItem() (ListItem, error) {
	var item ListItem
	if len(l.visibleItems) == 0 {
		return item, errors.New("no items")
	}
	return l.visibleItems[l.cursor], nil
}

func (l List) IsFirstIndexFocused() bool {
	return l.pageStartIdx+l.cursor == 0
}

func (l List) IsLastIndexFocused() bool {
	if len(l.items) == 0 {
		return true
	}
	return l.pageStartIdx+l.cursor == len(l.items)-1
}

func (l List) updateVisibleItems() []ListItem {
	start := l.pageStartIdx
	end := min(start+l.pageSize(), len(l.items))
	return l.items[start:end]
}

func (l List) nextPageStartIdx(offset int) int {
	start := l.pageStartIdx + offset
	if start+l.pageSize() > len(l.items) || start < 0 {
		return l.pageStartIdx
	}
	return start
}

func (l List) pageSize() int {
	return l.height
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
