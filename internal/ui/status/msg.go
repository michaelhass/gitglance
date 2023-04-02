package status

import "github.com/michaelhass/gitglance/internal/git"

// Git status

type initializedMsg struct {
	statusMsg statusUpdateMsg
	diffMsg   loadedDiffMsg
}

type statusUpdateMsg struct {
	err            error
	workTreeStatus git.WorkTreeStatus
}

type loadedDiffMsg struct {
	err  error
	diff string
}

// File handling

type focusItemMsg struct {
	item FileListItem
}

type selectItemMsg struct {
	item FileListItem
}
