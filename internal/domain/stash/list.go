package stash

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/list"
)

type ListItem struct {
	entry git.StashEntry
}

func (item ListItem) Render() string {
	return item.entry.Message()
}

type ListModel struct {
	listModel list.Model
}

func DefaultListItemHandler() list.ItemHandler {
	return func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case list.SelectItemMsg:
			if item, ok := msg.Item.(ListItem); ok {
				return popEntry(item.entry)
			}
			return nil
		case list.DeleteItemMsg:
			if item, ok := msg.Item.(ListItem); ok {
				return dropEntry(item.entry)
			}
			return nil
		case list.CustomItemMsg:
			if item, ok := msg.Item.(ListItem); ok {
				return dropEntry(item.entry)
			}
		}
		return nil
	}
}

func DefaultKeyMap() list.KeyMap {
	keyMap := list.NewKeyMap("", "pop", "drop")
	keyMap.All.SetEnabled(false)
	keyMap.Edit.SetEnabled(false)
	keyMap.CustomKeys = []key.Binding{
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "apply")),
	}
	return keyMap
}

func NewListModel(title string, keyMap list.KeyMap, itemHandler list.ItemHandler) ListModel {
	listModel := list.New(
		title,
		itemHandler,
		keyMap,
	)
	return ListModel{listModel: listModel}
}

func (sl ListModel) Init() tea.Cmd {
	return Load
}

func (sl ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(LoadedMsg); ok {
		var items []list.Item
		for _, entry := range msg.Stash {
			items = append(items, ListItem{entry: entry})
		}
		listModel, cmd := sl.listModel.SetItems(items)
		sl.listModel = listModel
		cmds = append(cmds, cmd)
	}

	listModel, cmd := sl.listModel.Update(msg)
	sl.listModel = listModel
	cmds = append(cmds, cmd)

	return sl, tea.Batch(cmds...)
}

func (sl ListModel) View() string {
	return sl.listModel.View()
}

func (sl ListModel) SetSize(width, height int) ListModel {
	sl.listModel = sl.listModel.SetSize(width, height)
	return sl
}

func (sl ListModel) Title() string {
	return sl.listModel.Title()
}

type ListDialogContent struct {
	ListModel
	width, height int
}
