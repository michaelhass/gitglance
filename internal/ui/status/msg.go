package status

import "github.com/michaelhass/gitglance/internal/git"

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
