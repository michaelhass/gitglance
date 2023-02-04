package status

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	styles "github.com/michaelhass/gitglance/internal/ui/Styles"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

var (
	itemStyle         = styles.TextSyle.Copy()
	focusedItemStyle  = styles.FocusTextStyle.Copy()
	inactiveItemStyle = styles.InactiveTextStyle.Copy()
)

type FileListItem struct {
	path, accessory string
}

func (item FileListItem) String() string {
	if len(item.accessory) == 0 {
		return item.path
	}
	return fmt.Sprintf("%s %s", item.accessory, item.path)
}

func NewFileListItem(path, accessory string) FileListItem {
	return FileListItem{
		path:      path,
		accessory: accessory,
	}
}

type FileList struct {
	items        []FileListItem
	visibleItems []FileListItem
	title        string
	width        int
	height       int
	cursor       int
	pageStartIdx int
	isFocused    bool
}

func NewFileList(title string) FileList {
	return FileList{title: title}
}
func (l FileList) Init() tea.Cmd {
	return nil
}

func (l FileList) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if l.cursor > 0 {
				l.cursor -= 1
				break
			}
			l.pageStartIdx = l.nextPageStartIdx(-1)
			l.visibleItems = l.updateVisibleItems()
		case "down":
			if l.cursor < len(l.visibleItems)-1 {
				l.cursor += 1
				break
			}
			l.pageStartIdx = l.nextPageStartIdx(1)
			l.visibleItems = l.updateVisibleItems()
		}
	}
	return l, nil
}

func (l FileList) View() string {
	return l.rendered()
}

func (l FileList) rendered() string {
	var renderedItems = make([]string, len(l.visibleItems))

	for i, item := range l.visibleItems {
		style := itemStyle
		if !l.isFocused {
			style = inactiveItemStyle
		} else if i == l.cursor {
			style = focusedItemStyle
		}

		itemString := item.String()
		itemString = itemString[:min(len(itemString), l.width-1)]

		renderedItems[i] = style.Width(l.width).Render(itemString)
	}

	return lipgloss.
		NewStyle().
		Height(l.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedItems...))
}

func (l FileList) Title() string {
	return l.title
}

func (l FileList) SetSize(width, height int) container.Content {
	l.width = width
	l.height = height
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l FileList) SetIsFocused(isFocused bool) container.Content {
	l.isFocused = isFocused
	return l
}

func (l FileList) SetFileListItems(items []FileListItem) FileList {
	l.items = items
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l FileList) updateVisibleItems() []FileListItem {
	start := l.pageStartIdx
	end := min(start+l.pageSize(), len(l.items))
	return l.items[start:end]
}

func (l FileList) nextPageStartIdx(offset int) int {
	start := l.pageStartIdx + offset
	if start+l.pageSize() > len(l.items) || start < 0 {
		return l.pageStartIdx
	}
	return start
}

func (l FileList) pageSize() int {
	return l.height
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
