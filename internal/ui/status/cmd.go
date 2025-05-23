package status

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/editor"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/list"
	"github.com/michaelhass/gitglance/internal/ui/refresh"
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

func refreshStatus() tea.Cmd {
	return tea.Sequence(
		updateWorkTreeStatus,
		list.ForceFocusUpdate,
	)
}

func stageFile(path string) tea.Cmd {
	return tea.Sequence(
		workTreeUpdateWithCmd(func() error {
			return git.StageFile(path)
		}),
		list.ForceFocusUpdate,
	)
}

func stageAll() tea.Cmd {
	return tea.Sequence(
		workTreeUpdateWithCmd(func() error {
			return git.StageAll()
		}),
		list.ForceFocusUpdate,
	)
}

func unstageFile(path string) tea.Cmd {
	return tea.Sequence(
		workTreeUpdateWithCmd(func() error {
			return git.UnstageFile(path)
		}),
		list.ForceFocusUpdate,
	)
}

func unstageAll() tea.Cmd {
	return tea.Sequence(
		workTreeUpdateWithCmd(func() error {
			return git.UnstageAll()
		}),
		list.ForceFocusUpdate,
	)
}

func deleteFile(path string, isUntracked bool) tea.Cmd {
	return tea.Sequence(
		workTreeUpdateWithCmd(func() error {
			return git.ResetFile(path, isUntracked)
		}),
		list.ForceFocusUpdate,
	)
}

func openFile(path string) tea.Cmd {
	return tea.ExecProcess(
		editor.OpenFileCmdDefault(
			path,
			editor.WithCmdString(git.CoreEditorValue),
			editor.WithCmdString(git.CoreGlobalEditorValue),
		),
		func(err error) tea.Msg {
			return refresh.Msg{}
		},
	)
}

func workTreeUpdateWithCmd(cmdFunc func() error) tea.Cmd {
	return func() tea.Msg {
		if err := cmdFunc(); err != nil {
			return statusUpdateMsg{Err: err}
		}

		return updateWorkTreeStatus()
	}
}

func updateWorkTreeStatus() tea.Msg {
	var (
		workTreeStatus git.WorkTreeStatus
		msg            statusUpdateMsg
		err            error
	)
	workTreeStatus, err = git.Status()
	if err != nil {
		msg.Err = err
		return msg
	}
	msg.WorkTreeStatus = workTreeStatus
	return msg
}

type loadedDiffMsg struct {
	Err  error
	Diff string
}

func showEmptyDiff() tea.Msg {
	return loadedDiffMsg{}
}

func diffFile(opt git.DiffOptions) func() tea.Msg {
	return func() tea.Msg {
		var (
			msg  loadedDiffMsg
			err  error
			diff string
		)

		msg.Diff = opt.FilePath
		diff, err = git.Diff(opt)
		if err != nil {
			msg.Err = err
			return msg
		}
		msg.Diff = diff

		return msg
	}
}
