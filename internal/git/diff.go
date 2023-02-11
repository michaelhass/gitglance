package git

import "os/exec"

type DiffOption struct {
	FilePath         string
	IsStaged         bool
	IsNameStatusOnly bool
	IsUntracked      bool
}

var untrackedFileDiffArgs = [3]string{
	"--no-index",
	"--",
	"/dev/null",
}

func newDiffCmd(opt DiffOption) exec.Cmd {
	args := []string{"diff"}

	if opt.IsStaged {
		args = append(args, "--cached")
	}

	if opt.IsNameStatusOnly {
		args = append(args, "--name-status")
	}

	if opt.IsUntracked {
		args = append(args, untrackedFileDiffArgs[:3]...)
	}

	if len(opt.FilePath) > 0 {
		args = append(args, opt.FilePath)
	}

	return *exec.Command("git", args...)
}
