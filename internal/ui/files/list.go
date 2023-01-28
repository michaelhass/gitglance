package files

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	styles "github.com/michaelhass/gitglance/internal/ui/Styles"
)

var (
	itemStyle        = lipgloss.NewStyle()
	focusedItemStyle = itemStyle.Foreground(lipgloss.Color("170"))
	disabledStyle    = lipgloss.NewStyle()
)

type ListItem struct {
	path, accessory string
}

func (item ListItem) String() string {
	if len(item.accessory) == 0 {
		return item.path
	}
	return fmt.Sprintf("%s %s", item.accessory, item.path)
}

func NewListItem(path, accessory string) ListItem {
	return ListItem{
		path:      path,
		accessory: accessory,
	}
}

type List struct {
	items         []ListItem
	title         string
	width, height int
	cursor        int
	isEnabled     bool
}

func NewList(title string) List {
	return List{title: title}
}
func (l List) Init() tea.Cmd {
	return nil
}

func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
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

func (l List) View() string {
	return l.rendered()
}

func (l List) rendered() string {
	var builder strings.Builder

	builder.WriteString(styles.TitleStyle.Render(l.title))
	builder.WriteString("\n")
	builder.WriteString("\n")

	if len(l.items) == 0 {
		return builder.String()
	}

	for i, item := range l.items {
		style := itemStyle
		if !l.isEnabled {
			style = disabledStyle
		} else if i == l.cursor {
			style = focusedItemStyle
		}

		builder.WriteString(style.Copy().Width(l.width - 2).Render(item.String()))
		builder.WriteString("\n")
	}

	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Height(l.height).Render(builder.String())
}

func (l *List) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *List) SetItems(items []ListItem) {
	l.items = items
}

func (l *List) SetIsEnabled(isEnabled bool) {
	l.isEnabled = isEnabled
}

func (l *List) IsEnabled() bool {
	return l.isEnabled
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
