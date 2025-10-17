package stash

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/list"
)

type StashListItem struct {
	entry git.StashEntry
}

func (item StashListItem) Render() string {
	return item.entry.Message()
}

type StashList struct {
	listModel list.Model
}

func NewStashList(title string, keyMap list.KeyMap) StashList {
	listModel := list.New(title, func(msg tea.Msg) tea.Cmd { return nil }, keyMap)
	return StashList{listModel: listModel}
}

func (sl StashList) Init() tea.Cmd {
	return Load
}

func (sl StashList) Update(msg tea.Msg) (StashList, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(LoadedMsg); ok {
		var items []list.Item
		for _, entry := range msg.Stash {
			items = append(items, StashListItem{entry: entry})
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

func (sl StashList) View() string {
	return sl.listModel.View()
}

func (sl StashList) SetSize(width, height int) StashList {
	sl.listModel = sl.listModel.SetSize(width, height)
	return sl
}

func (sl StashList) Title() string {
	return sl.listModel.Title()
}

type ApplyDialogContent struct {
	StashList
	width, height int
}
