package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
)

type mockWorktree struct {
	status git.Status
}

func (wt mockWorktree) Status() (git.Status, error) {
	return wt.status, nil
}

func TestStatus(t *testing.T) {

	const (
		unmodifiedPath   = "unmodified"
		untrackedPath    = "untracked"
		onlyStagedPath   = "only_staged"
		onlyUnstagedPath = "only_unstaged"
		bothPath         = "both"
	)

	sourceStatus := git.Status{}
	sourceStatus[unmodifiedPath] = &git.FileStatus{Staging: git.Unmodified, Worktree: git.Unmodified}
	sourceStatus[untrackedPath] = &git.FileStatus{Staging: git.Untracked, Worktree: git.Untracked}
	sourceStatus[onlyStagedPath] = &git.FileStatus{Staging: git.Modified, Worktree: git.Unmodified}
	sourceStatus[onlyUnstagedPath] = &git.FileStatus{Staging: git.Unmodified, Worktree: git.Modified}
	sourceStatus[bothPath] = &git.FileStatus{Staging: git.Modified, Worktree: git.Modified}

	wt := newWorkTree(mockWorktree{status: sourceStatus})

	status, err := wt.Status()
	if err != nil {
		t.Errorf("Failed reading status. Error: %v", err)
	}

	if status.contains(status.Staged, unmodifiedPath, Unmodified) ||
		status.contains(status.Unstaged, unmodifiedPath, Unmodified) {
		t.Errorf("Failed to remove unmodified status.")
	}

	if status.contains(status.Staged, untrackedPath, Untracked) ||
		!status.contains(status.Unstaged, untrackedPath, Untracked) {
		t.Errorf("Failed to remove unmodified status.")
	}

	if !status.contains(status.Staged, onlyStagedPath, Modified) ||
		status.contains(status.Unstaged, onlyStagedPath, Modified) {
		t.Errorf("File should be included in staged but not in unstaged.")
	}

	if status.contains(status.Staged, onlyUnstagedPath, Modified) ||
		!status.contains(status.Unstaged, onlyUnstagedPath, Modified) {
		t.Errorf("File should be included in unstaged but not in staged.")
	}

	if !status.contains(status.Staged, bothPath, Modified) ||
		!status.contains(status.Unstaged, bothPath, Modified) {
		t.Errorf("File should be included in both staged and unstaged.")
	}
}
