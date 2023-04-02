package status

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

func initializeStatus() func() tea.Msg {
	return func() tea.Msg {
		var (
			msg            initializedMsg
			workTreeStatus git.WorkTreeStatus
			unstagedFiles  git.FileStatusList
			isUntracked    bool
			err            error
		)

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.statusMsg.err = err
			return msg
		}
		msg.statusMsg.workTreeStatus = workTreeStatus

		unstagedFiles = msg.statusMsg.workTreeStatus.Unstaged
		if len(unstagedFiles) == 0 {
			return msg
		}

		isUntracked = unstagedFiles[0].Code == git.Untracked
		diffMsg, ok := diff(
			git.DiffOption{
				FilePath:    unstagedFiles[0].Path,
				IsUntracked: isUntracked,
			},
		)().(loadedDiffMsg)
		if !ok {
			diffMsg.err = errors.New("unable to load diff")
			return msg
		}

		msg.diffMsg = diffMsg
		return msg
	}
}

func stageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
			err            error
		)

		err = git.StageFile(path)
		if err != nil {
			msg.err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.workTreeStatus = workTreeStatus

		return msg
	}
}

func unstageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
			err            error
		)

		err = git.UnstageFile(path)
		if err != nil {
			msg.err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.workTreeStatus = workTreeStatus

		return msg
	}
}

func diff(opt git.DiffOption) func() tea.Msg {
	return func() tea.Msg {
		var (
			msg  loadedDiffMsg
			err  error
			diff string
		)

		diff, err = git.Diff(opt)
		if err != nil {
			msg.err = err
			return msg
		}
		msg.diff = diff

		return msg
	}
}
