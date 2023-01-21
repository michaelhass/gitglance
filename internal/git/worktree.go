package git

import (
	"github.com/go-git/go-git/v5"
)

type worktreeBridge interface {
	Status() (git.Status, error)
}

type Worktree struct {
	wtb worktreeBridge
}

func newWorkTree(wtb worktreeBridge) *Worktree {
	return &Worktree{wtb: wtb}
}

func (wt *Worktree) Status() (Status, error) {
	var status Status
	bridgeStatus, err := wt.wtb.Status()

	if err != nil {
		return status, err
	}

	for path, fileStatus := range bridgeStatus {
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
