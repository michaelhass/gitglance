package git

import (
	"github.com/go-git/go-git/v5"
)

type Worktree struct {
	wt *git.Worktree
}

func newWorkTree(wt *git.Worktree) *Worktree {
	return &Worktree{wt: wt}
}

func (wt *Worktree) Status() (Status, error) {
	return wt.readStatus(func() (git.Status, error) {
		return wt.wt.Status()
	})
}

func (wt *Worktree) readStatus(readStatus func() (git.Status, error)) (Status, error) {
	var status Status

	srcStatus, err := readStatus()
	if err != nil {
		return status, err
	}

	for path, fileStatus := range srcStatus {
		if code := StatusCode(fileStatus.Worktree); code != Unmodified {
			status.Unstaged = append(
				status.Unstaged,
				FileStatus{
					Path: path,
					Code: code,
				},
			)
		}
		if code := StatusCode(fileStatus.Staging); code != Unmodified && code != Untracked {
			status.Staged = append(
				status.Staged,
				FileStatus{
					Path: path,
					Code: code,
				},
			)
		}
	}
	return status, nil
}

func (wt *Worktree) StageFile(path string) error {
	_, err := wt.wt.Add(path)
	return err
}
