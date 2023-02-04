package status

import (
	"fmt"
	"strings"

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
	items     []FileListItem
	title     string
	width     int
	height    int
	cursor    int
	isFocused bool
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
			l.cursor = max(l.cursor-1, 0)
		case "down":
			l.cursor = min(l.cursor+1, len(l.items)-1)
		}
	}
	return l, nil
}

func (l FileList) View() string {
	return l.rendered()
}

func (l FileList) rendered() string {
	var builder strings.Builder

	for i, item := range l.items {
		style := itemStyle
		if !l.isFocused {
			style = inactiveItemStyle
		} else if i == l.cursor {
			style = focusedItemStyle
		}

		builder.WriteString(style.Copy().Width(l.width).Render(item.String()))
		builder.WriteString("\n")
	}

	return lipgloss.NewStyle().Height(l.height).Render(builder.String())
}

func (l FileList) Title() string {
	return l.title
}

func (l FileList) SetSize(width, height int) container.Content {
	l.width = width
	l.height = height
	return l
}

func (l FileList) SetIsFocused(isFocused bool) container.Content {
	l.isFocused = isFocused
	return l
}

func (l FileList) SetFileListItems(items []FileListItem) FileList {
	l.items = items
	return l
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
