package cmd

import "github.com/michaelhass/gitglance/internal/git"

// Git status

type StatusUpdateMsg struct {
	Err            error
	WorkTreeStatus git.WorkTreeStatus
}

type LoadedDiffMsg struct {
	Err  error
	Diff string
}
