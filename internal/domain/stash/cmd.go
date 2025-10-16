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
