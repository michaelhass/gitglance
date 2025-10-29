package stash

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/err"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog/confirm"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog/info"
)

type EntryCmdType byte

const (
	CreatedEntryCmdType EntryCmdType = iota
	AppliedEntryCmdType
	PoppedEntryCmdType
	DroppedEntryCmdType
)

type EntryCmdExecuted struct {
	CmdType EntryCmdType
	Entry   git.StashEntry
	err     error
}

func (ece EntryCmdExecuted) Err() error {
	return ece.err
}

func (ece EntryCmdExecuted) ErrorTitle() string {
	switch ece.CmdType {
	case CreatedEntryCmdType:
		return "Create stash error"
	case AppliedEntryCmdType:
		return "Apply stash entry error"
	case PoppedEntryCmdType:
		return "Pop stash entry error"
	case DroppedEntryCmdType:
		return "Drop stash entry error"
	}
	return ""
}

func (ece EntryCmdExecuted) ErrorDescription() string {
	return ece.err.Error()
}

func executionErrHandler(msg tea.Msg) tea.Cmd {
	if errMsg, ok := msg.(err.Msg); ok && errMsg.Err() != nil {
		errDialogContent := info.NewDialogContentWithErrMsg(errMsg)
		return dialog.Show(errDialogContent, nil, dialog.CenterDisplayMode)
	}
	return nil
}

func CreateWithUntracked(msg string) tea.Cmd {
	return func() tea.Msg {
		opts := git.CreateStashOpts{}
		opts.WithUntracked = true
		opts.Message = msg
		err := git.CreateStash(opts)
		return EntryCmdExecuted{CmdType: CreatedEntryCmdType, err: err}
	}
}
func ShowCreateWithUntrackedConfirmation(onClose tea.Cmd) tea.Cmd {
	confirmModel := confirm.
		New("Stash", "Do you want to stash all changes?").
		WithTextInput(
			"Message...",
			func(message string) tea.Cmd {
				return CreateWithUntracked(message)
			},
		)
	confirmDialog := confirm.
		NewDialogContent(confirmModel).
		WithErrHandler(executionErrHandler)
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

func showActionConfirmation(confirmCmd tea.Cmd, msg string) tea.Cmd {
	dc := confirm.NewDialogContent(confirm.New("Stash", msg).
		WithOnConfirmCmd(confirmCmd)).
		WithErrHandler(executionErrHandler)
	return dialog.Show(dc, Load, dialog.CenterDisplayMode)
}

func applyEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.ApplyStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: AppliedEntryCmdType,
			Entry:   entry,
			err:     err,
		}
	}
}

func popEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.PopStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: PoppedEntryCmdType,
			Entry:   entry,
			err:     err,
		}
	}
}

func dropEntry(entry git.StashEntry) tea.Cmd {
	return func() tea.Msg {
		err := git.DropStashEntry(entry)
		return EntryCmdExecuted{
			CmdType: DroppedEntryCmdType,
			Entry:   entry,
			err:     err,
		}
	}
}
