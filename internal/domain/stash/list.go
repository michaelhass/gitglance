package stash

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/components/list"
	"github.com/michaelhass/gitglance/internal/core/ui/style"
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

func NewStashList() StashList {
	keyMap := list.NewKeyMap("", "Apply stash", "")
	keyMap.All.SetEnabled(false)
	keyMap.Edit.SetEnabled(false)
	keyMap.Delete.SetEnabled(false)
	listModel := list.New("Stash", func(msg tea.Msg) tea.Cmd { return nil }, keyMap)

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
	return lipgloss.JoinVertical(
		lipgloss.Left,
		style.Title.Render(sl.listModel.Title()),
		"",
		sl.listModel.View(),
	)
}

func (sl StashList) SetSize(width, height int) StashList {
	sl.listModel = sl.listModel.SetSize(width, height)
	return sl
}

type DialogContent struct {
	StashList
}

func NewDialogConent(stashList StashList) DialogContent {
	stashList.listModel, _ = stashList.listModel.UpdateFocus(true)
	return DialogContent{StashList: stashList}
}

func (dc DialogContent) Init() tea.Cmd {
	return dc.StashList.Init()
}

func (dc DialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	model, cmd := dc.StashList.Update(msg)
	dc.StashList = model
	return dc, cmd
}

func (dc DialogContent) View() string {
	return dc.StashList.View()
}

func (dc DialogContent) SetSize(width, height int) dialog.Content {
	dc.StashList = dc.StashList.SetSize(width, height)
	return dc
}

func (dc DialogContent) Help() []key.Binding {
	return dc.listModel.KeyMap().ShortHelp()
}
