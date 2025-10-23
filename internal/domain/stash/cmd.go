package stash

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog/confirm"
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

func ShowListDialog(onClose tea.Cmd) tea.Cmd {
	stashList := NewListModel("Stash", DefaultKeyMap(), DefaultListItemHandler())
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
		err := git.PopStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: PoppedEntryCmdType,
			Entry:   entry,
			Err:     err,
		}
	}
}

func dropEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.DropStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: DroppedEntryCmdType,
			Entry:   entry,
			Err:     err,
		}
	}
}
