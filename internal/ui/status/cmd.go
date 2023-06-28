package status

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

type InitializedMsg struct {
	StatusMsg statusUpdateMsg
	DiffMsg   loadedDiffMsg
}

func initializeStatus() tea.Cmd {
	return func() tea.Msg {
		var (
			msg            InitializedMsg
			workTreeStatus git.WorkTreeStatus
			unstagedFiles  git.FileStatusList
			isUntracked    bool
			err            error
		)

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.StatusMsg.Err = err
			return msg
		}
		msg.StatusMsg.WorkTreeStatus = workTreeStatus

		unstagedFiles = msg.StatusMsg.WorkTreeStatus.UnstagedFiles()
		if len(unstagedFiles) == 0 {
			return msg
		}

		isUntracked = unstagedFiles[0].IsUntracked()
		diffMsg, ok := diffFile(
			git.DiffOption{
				FilePath:    unstagedFiles[0].Path,
				IsUntracked: isUntracked,
			},
		)().(loadedDiffMsg)
		if !ok {
			diffMsg.Err = errors.New("unable to load diff")
			return msg
		}

		msg.DiffMsg = diffMsg
		return msg
	}
}

type focusSectionMsg struct {
	section section
}

func focusSection(section section) tea.Cmd {
	return func() tea.Msg {
		return focusSectionMsg{
			section: section,
		}
	}
}

type statusUpdateMsg struct {
	Err            error
	WorkTreeStatus git.WorkTreeStatus
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

func unstageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
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

type loadedDiffMsg struct {
	Err  error
	Diff string
}

func diffFile(opt git.DiffOption) func() tea.Msg {
	return func() tea.Msg {
		var (
			msg  loadedDiffMsg
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
