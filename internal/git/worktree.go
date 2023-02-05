package git

import (
	"os/exec"
	"strings"
)

type Worktree struct{}

func newWorkTree() *Worktree {
	return &Worktree{}
}

func (wt *Worktree) Status() (Status, error) {
	var status Status

	unstagedFiles, err := wt.fileStatusListFromDiff(DiffOption{IsNameStatusOnly: true})
	if err != nil {
		return status, err
	}

	untrackedFilesPaths, err := wt.untrackedFiles()
	if err != nil {
		return status, nil
	}

	untrackedFiles := fileStatusListForUntrackedFiles(untrackedFilesPaths)
	unstagedFiles = append(unstagedFiles, untrackedFiles...)

	stagedFiles, err := wt.fileStatusListFromDiff(DiffOption{IsNameStatusOnly: true, IsStaged: true})
	if err != nil {
		return status, err
	}

	status = newStatus(unstagedFiles, stagedFiles)
	return status, nil
}

func (wt *Worktree) fileStatusListFromDiff(opt DiffOption) (FileStatusList, error) {
	diff, err := wt.Diff(opt)
	if err != nil {
		return nil, err
	}
	return fileStatusListFromDiffString(diff)
}

func (wt *Worktree) StageFile(path string) error {
	cmd := exec.Command("git", "add", path)
	return cmd.Run()
}

func (wt *Worktree) UnstageFile(path string) error {
	cmd := exec.Command("git", "restore", "--staged", path)
	return cmd.Run()
}

func (wt *Worktree) untrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}
		files = append(files, line)
	}
	return files, nil

}

type DiffOption struct {
	FilePath         string
	IsStaged         bool
	IsNameStatusOnly bool
}

func (wt *Worktree) Diff(opt DiffOption) (string, error) {
	cmd := newDiffCmd(opt)
	out, err := cmd.Output()
	return string(out), err
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
