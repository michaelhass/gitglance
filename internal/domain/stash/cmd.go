package stash

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog/confirm"
	"github.com/michaelhass/gitglance/internal/core/ui/components/list"
)

type CreatedMsg struct {
	err error
}

func CreateWithUntracked() tea.Msg {
	err := git.StashAll()
	return CreatedMsg{err: err}
}

func ShowCreateWithUntrackedConfirmation(onClose tea.Cmd) tea.Cmd {
	confirmDialog := confirm.NewDialogConent(
		confirm.New(
			"Stash", "Do you want to stash all changes?",
			CreateWithUntracked,
		),
	)
	return dialog.Show(confirmDialog, onClose, dialog.CenterDisplayMode)
}

type LoadedMsg struct {
	Stash git.Stash
	Err   error
}

func Load() tea.Msg {
	stash, err := git.GetStash()
	return LoadedMsg{Stash: stash, Err: err}
}

func ShowApplyDialog(onClose tea.Cmd) tea.Cmd {
	keyMap := list.NewKeyMap("", "apply stash", "")
	keyMap.All.SetEnabled(false)
	keyMap.Edit.SetEnabled(false)
	keyMap.Delete.SetEnabled(false)
	stashList := NewStashList("Apply stash", keyMap)
	return dialog.Show(NewApplyDialogConent(stashList), onClose, dialog.CenterDisplayMode)
}
