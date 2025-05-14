package filelist

import (
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/ui/style"
)

var (
	itemStyle         = style.Text.Copy()
	focusedItemStyle  = style.FocusText.Copy()
	inactiveItemStyle = style.InactiveText.Copy()
)

type ItemHandler func(msg tea.Msg) tea.Cmd

type Model struct {
	items         []Item
	visibleItems  []Item
	itemHandler   ItemHandler
	keys          KeyMap
	title         string
	width         int
	height        int
	cursor        int
	pageStartIdx  int
	isFocused     bool
	lastFocuedIdx int
}

func New(title string, itemHandler ItemHandler, keys KeyMap) Model {
	return Model{title: title, itemHandler: itemHandler, keys: keys}
}
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.isFocused {
		return m, nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case forceFocusUpdateMsg:
		model, cmd := m.updateFocus(m.isFocused, true)
		m = model
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor == 0 || len(m.visibleItems) == 0 {
				if m.IsFirstIndexFocused() {
					cmds = append(cmds, m.itemHandler(TopNoMoreFocusableItems{}))
				} else {
					m.pageStartIdx = m.nextPageStartIdx(-1)
					m.visibleItems = m.updateVisibleItems()
				}
				break
			}
			m.cursor -= 1
			cmd := m.itemHandler(FocusItemMsg{Item: m.visibleItems[m.cursor]})
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Down):
			if m.cursor >= len(m.visibleItems)-1 {
				if m.IsLastIndexFocused() {
					cmds = append(cmds, m.itemHandler(BottomNoMoreFocusableItems{}))
				} else {
					m.pageStartIdx = m.nextPageStartIdx(1)
					m.visibleItems = m.updateVisibleItems()
				}
				break
			}
			m.cursor += 1
			cmd := m.itemHandler(FocusItemMsg{Item: m.visibleItems[m.cursor]})
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Enter):
			if len(m.visibleItems) == 0 {
				break
			}
			item := m.visibleItems[m.cursor]
			cmd := m.itemHandler(SelectItemMsg{Item: item})
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.All):
			cmd := m.itemHandler(SelectAllItemMsg{Items: m.items})
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Delete):
			if len(m.visibleItems) == 0 {
				break
			}
			item := m.visibleItems[m.cursor]
			cmd := m.itemHandler(DeleteItemMsg{Item: item})
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) UpdateFocus(isFocused bool) (Model, tea.Cmd) {
	return m.updateFocus(isFocused, false)
}

func (m Model) updateFocus(isFocused bool, isForced bool) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.isFocused = isFocused
	isAlreadyFocused := m.lastFocuedIdx == m.cursor
	if !isFocused {
		m.lastFocuedIdx = -1
	} else if (!isAlreadyFocused || isForced) && isFocused && len(m.visibleItems) > 0 {
		m.lastFocuedIdx = m.cursor
		cmd = m.itemHandler(FocusItemMsg{Item: m.visibleItems[m.cursor]})
	} else if isFocused && len(m.items) == 0 {
		cmd = m.itemHandler(NoItemsMsg{})
	}
	return m, cmd
}

func (m Model) View() string {
	var renderedItems = make([]string, len(m.visibleItems))

	for i, item := range m.visibleItems {
		style := itemStyle
		if !m.isFocused {
			style = inactiveItemStyle
		} else if i == m.cursor {
			style = focusedItemStyle
		}
		renderedItems[i] = style.MaxHeight(1).MaxWidth(m.width - 1).Render(item.String())
	}

	return lipgloss.
		NewStyle().
		Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedItems...))
}

func (m Model) Title() string {
	return m.title
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	m.visibleItems = m.updateVisibleItems()
	return m
}

func (m Model) KeyMap() help.KeyMap {
	return m.keys
}

func (m Model) SetItems(items []Item) (Model, tea.Cmd) {
	m.items = items
	m.visibleItems = m.updateVisibleItems()

	// Check out of bounds due to content change
	if len(m.visibleItems) > 0 && m.cursor >= len(m.visibleItems) {
		m.cursor -= 1
		m.pageStartIdx = m.nextPageStartIdx(-1)
		m.visibleItems = m.updateVisibleItems()
	}

	if !m.isFocused || len(m.visibleItems) == 0 {
		return m, nil
	}

	return m, m.itemHandler(FocusItemMsg{Item: m.visibleItems[m.cursor]})
}

func (m Model) FocusedItem() (Item, error) {
	var item Item
	if len(m.visibleItems) == 0 {
		return item, errors.New("no items")
	}
	return m.visibleItems[m.cursor], nil
}

func (m Model) IsFirstIndexFocused() bool {
	if len(m.items) == 0 {
		return true
	}
	return m.pageStartIdx+m.cursor == 0
}

func (m Model) IsLastIndexFocused() bool {
	if len(m.items) == 0 {
		return true
	}
	return m.pageStartIdx+m.cursor == len(m.items)-1
}

func (m Model) updateVisibleItems() []Item {
	if len(m.items) == 0 || m.height <= 0 {
		return []Item{}
	}
	start := m.pageStartIdx
	end := min(start+m.pageSize(), len(m.items))
	return m.items[start:end]
}

func (m Model) nextPageStartIdx(offset int) int {
	start := m.pageStartIdx + offset
	if start+m.pageSize() > len(m.items) || start < 0 {
		return m.pageStartIdx
	}
	return start
}

func (m Model) pageSize() int {
	return m.height
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
