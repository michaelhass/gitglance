package status

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

type initializedMsg struct {
	statusMsg statusUpdateMsg
	diffMsg   loadedDiffMsg
}

func initializeStatus() tea.Cmd {
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
			msg.statusMsg.Err = err
			return msg
		}
		msg.statusMsg.WorkTreeStatus = workTreeStatus

		unstagedFiles = msg.statusMsg.WorkTreeStatus.UnstagedFiles()
		if len(unstagedFiles) == 0 {
			return msg
		}

		isUntracked = unstagedFiles[0].IsUntracked()
		diffMsg, ok := diffFile(
			git.DiffOptions{
				FilePath:    unstagedFiles[0].Path,
				IsUntracked: isUntracked,
			},
		)().(loadedDiffMsg)
		if !ok {
			diffMsg.Err = errors.New("unable to load diff")
			return msg
		}

		msg.diffMsg = diffMsg
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
	return workTreeUpdateWithCmd(func() error {
		return git.StageFile(path)
	})
}

func stageAll() func() tea.Msg {
	return workTreeUpdateWithCmd(func() error {
		return git.StageAll()
	})
}

func unstageFile(path string) func() tea.Msg {
	return workTreeUpdateWithCmd(func() error {
		return git.UnstageFile(path)
	})
}

func unstageAll() func() tea.Msg {
	return workTreeUpdateWithCmd(func() error {
		return git.UnstageAll()
	})
}

func deleteFile(path string, isUntracked bool) func() tea.Msg {
	return workTreeUpdateWithCmd(func() error {
		return git.ResetFile(path, isUntracked)
	})
}

func workTreeUpdateWithCmd(cmdFunc func() error) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
			err            error
		)

		err = cmdFunc()
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

func diffFile(opt git.DiffOptions) func() tea.Msg {
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
