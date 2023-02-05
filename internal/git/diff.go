package git

import "os/exec"

type DiffOption struct {
	FilePath         string
	IsStaged         bool
	IsNameStatusOnly bool
}

func newDiffCmd(opt DiffOption) exec.Cmd {
	args := []string{"diff"}

	if opt.IsStaged {
		args = append(args, "--cached")
	}

	if opt.IsNameStatusOnly {
		args = append(args, "--name-status")
	}

	if len(opt.FilePath) > 0 {
		args = append(args, opt.FilePath)
	}

	return *exec.Command("git", args...)
}
