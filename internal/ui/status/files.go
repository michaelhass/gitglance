package status

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/styles"
)

var (
	itemStyle         = styles.TextSyle.Copy()
	focusedItemStyle  = styles.FocusTextStyle.Copy()
	inactiveItemStyle = styles.InactiveTextStyle.Copy()
)

type focusItemMsg struct {
	item FileListItem
}

type selectItemMsg struct {
	item FileListItem
}

type FileListItem struct {
	fileStatus git.FileStatus
	path       string
	accessory  string
}

func (item FileListItem) String() string {
	if len(item.accessory) == 0 {
		return item.path
	}
	return fmt.Sprintf("%s %s", item.accessory, item.path)
}

func NewFileListItem(fileStatus git.FileStatus) FileListItem {
	var (
		path, accessory string
	)
	path = fileStatus.Path
	if len(fileStatus.Extra) > 0 {
		path = fmt.Sprintf("%s → %s", path, fileStatus.Extra)
	}

	accessory = fmt.Sprintf("[%s]", string(fileStatus.Code))

	return FileListItem{
		fileStatus: fileStatus,
		path:       path,
		accessory:  accessory,
	}
}

type ItemHandler func(msg tea.Msg) tea.Cmd

type FileList struct {
	items        []FileListItem
	visibleItems []FileListItem
	itemHandler  ItemHandler
	title        string
	width        int
	height       int
	cursor       int
	pageStartIdx int
	isFocused    bool
}

func NewFileList(title string, itemHandler ItemHandler) FileList {
	return FileList{title: title, itemHandler: itemHandler}
}
func (l FileList) Init() tea.Cmd {
	return nil
}

func (l FileList) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	if !l.isFocused || len(l.items) == 0 {
		return l, nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if l.cursor == 0 {
				l.pageStartIdx = l.nextPageStartIdx(-1)
				l.visibleItems = l.updateVisibleItems()
				break
			}
			l.cursor -= 1
			cmd := l.itemHandler(focusItemMsg{item: l.visibleItems[l.cursor]})
			cmds = append(cmds, cmd)

		case tea.KeyDown:
			if l.cursor >= len(l.visibleItems)-1 {
				l.pageStartIdx = l.nextPageStartIdx(1)
				l.visibleItems = l.updateVisibleItems()
				break
			}
			l.cursor += 1
			cmd := l.itemHandler(focusItemMsg{item: l.visibleItems[l.cursor]})
			cmds = append(cmds, cmd)

		case tea.KeyEnter:
			item := l.visibleItems[l.cursor]
			cmd := l.itemHandler(selectItemMsg{item: item})
			cmds = append(cmds, cmd)
		}

	default:
		if len(l.visibleItems) == 0 || l.cursor < len(l.visibleItems) {
			break
		}
		l.cursor = len(l.visibleItems) - 1
	}

	return l, tea.Batch(cmds...)
}

func (l FileList) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	var cmd tea.Cmd
	if isFocused && !l.isFocused && len(l.items) > 0 {
		cmd = l.itemHandler(focusItemMsg{item: l.visibleItems[l.cursor]})
	}
	l.isFocused = isFocused
	return l, cmd
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

func (l FileList) Title() string {
	return l.title
}

func (l FileList) SetSize(width, height int) container.Content {
	l.width = width
	l.height = height
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l FileList) SetFileListItems(items []FileListItem) FileList {
	l.items = items
	l.visibleItems = l.updateVisibleItems()
	return l
}

func (l FileList) FocusedItem() (FileListItem, error) {
	var item FileListItem
	if len(l.visibleItems) == 0 {
		return item, errors.New("no items")
	}
	return l.visibleItems[l.cursor], nil
}

func (l FileList) IsLastIndexFocused() bool {
	return l.pageStartIdx+l.cursor == len(l.items)-1
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
