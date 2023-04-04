package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

func StageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            StatusUpdateMsg
			err            error
		)

		err = git.StageFile(path)
		if err != nil {
			msg.Err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.Err = err
			return msg
		}
		msg.WorkTreeStatus = workTreeStatus

		return msg
	}
}

func UnstageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            StatusUpdateMsg
			err            error
		)

		err = git.UnstageFile(path)
		if err != nil {
			msg.Err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.Err = err
			return msg
		}
		msg.WorkTreeStatus = workTreeStatus

		return msg
	}
}

func Diff(opt git.DiffOption) func() tea.Msg {
	return func() tea.Msg {
		var (
			msg  LoadedDiffMsg
			err  error
			diff string
		)

		diff, err = git.Diff(opt)
		if err != nil {
			msg.Err = err
			return msg
		}
		msg.Diff = diff

		return msg
	}
}
