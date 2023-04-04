package status

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
	uicmd "github.com/michaelhass/gitglance/internal/ui/cmd"
)

func initializeStatus() func() tea.Msg {
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

		unstagedFiles = msg.StatusMsg.WorkTreeStatus.Unstaged
		if len(unstagedFiles) == 0 {
			return msg
		}

		isUntracked = unstagedFiles[0].Code == git.Untracked
		diffMsg, ok := uicmd.Diff(
			git.DiffOption{
				FilePath:    unstagedFiles[0].Path,
				IsUntracked: isUntracked,
			},
		)().(uicmd.LoadedDiffMsg)
		if !ok {
			diffMsg.Err = errors.New("unable to load diff")
			return msg
		}

		msg.DiffMsg = diffMsg
		return msg
	}
}

type InitializedMsg struct {
	StatusMsg uicmd.StatusUpdateMsg
	DiffMsg   uicmd.LoadedDiffMsg
}
