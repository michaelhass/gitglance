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
	keyMap := list.NewKeyMap("", "pop", "drop")
	keyMap.All.SetEnabled(false)
	keyMap.Edit.SetEnabled(false)
	keyMap.Delete.SetEnabled(true)
	stashList := NewStashList("Apply stash", keyMap, DefaultListItemHandler())
	return dialog.Show(NewApplyDialogConent(stashList), onClose, dialog.CenterDisplayMode)
}

type EntryCmdType byte

const (
	AppliedEntryCmdType EntryCmdType = iota
	PoppedEntryCmdType  EntryCmdType = iota
	DroppedEntryCmdType EntryCmdType = iota
)

type EntryCmdExecuted struct {
	CmdType EntryCmdType
	Entry   git.StashEntry
	Err     error
}

func applyEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.ApplyStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: AppliedEntryCmdType,
			Entry:   entry,
			Err:     err,
		}
	}
}

func popEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.ApplyStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: PoppedEntryCmdType,
			Entry:   entry,
			Err:     err,
		}
	}
}

func dropEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.ApplyStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: DroppedEntryCmdType,
			Entry:   entry,
			Err:     err,
		}
	}
}
