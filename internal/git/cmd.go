package git

import (
	"errors"
	"os/exec"
)

type gitCommand struct {
	cmd *exec.Cmd
}

func newGitCommand(args ...string) *gitCommand {
	cmd := exec.Command("git", args...)
	return &gitCommand{cmd: cmd}
}

func (gc *gitCommand) run() error {
	if err := gc.cmd.Run(); err != nil && !isExitError(err) {
		return err
	}
	return nil
}

func (gc *gitCommand) output() (string, error) {
	out, err := gc.cmd.Output()
	if err != nil && !isExitError(err) {
		return "", err
	}
	return string(out), nil
}

func isExitError(err error) bool {
	var ee *exec.ExitError
	return errors.As(err, &ee)
}

type statusOptions struct {
	isPorcelain     bool
	isNULTerminated bool
	hasBranch       bool
	isShort         bool
}

func newStatusCmd(opts statusOptions) *gitCommand {
	args := []string{"status"}

	if opts.isPorcelain {
		args = append(args, "--porcelain")
	}

	if opts.isShort {
		args = append(args, "--short")
	}

	if opts.isNULTerminated {
		args = append(args, "-z")
	}

	if opts.hasBranch {
		args = append(args, "-b")
	}

	return newGitCommand(args...)
}

type DiffOptions struct {
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

func newDiffCmd(opts DiffOptions) *gitCommand {
	args := []string{"diff"}

	if opts.IsStaged {
		args = append(args, "--cached")
	}

	if opts.IsNameStatusOnly {
		args = append(args, "--name-status")
	}

	if opts.IsUntracked {
		args = append(args, untrackedFileDiffArgs[:3]...)
	}

	if len(opts.FilePath) > 0 {
		args = append(args, opts.FilePath)
	}

	return newGitCommand(args...)
}

func removeFileCmd(filePath string) *exec.Cmd {
	return exec.Command("rm", "-rf", filePath)
}
