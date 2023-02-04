package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestStatus(t *testing.T) {

	const (
		unmodifiedPath   = "unmodified"
		untrackedPath    = "untracked"
		onlyStagedPath   = "only_staged"
		onlyUnstagedPath = "only_unstaged"
		bothPath         = "both"
	)

	wt := newWorkTree(nil)

	status, err := wt.readStatus(func() (git.Status, error) {
		status := git.Status{}
		status[unmodifiedPath] = &git.FileStatus{Staging: git.Unmodified, Worktree: git.Unmodified}
		status[untrackedPath] = &git.FileStatus{Staging: git.Untracked, Worktree: git.Untracked}
		status[onlyStagedPath] = &git.FileStatus{Staging: git.Modified, Worktree: git.Unmodified}
		status[onlyUnstagedPath] = &git.FileStatus{Staging: git.Unmodified, Worktree: git.Modified}
		status[bothPath] = &git.FileStatus{Staging: git.Modified, Worktree: git.Modified}
		return status, nil
	})

	if err != nil {
		t.Errorf("Failed reading status. Error: %v", err)
	}

	if status.Staged.Contains(unmodifiedPath, Unmodified) ||
		status.Unstaged.Contains(unmodifiedPath, Unmodified) {
		t.Errorf("Failed to remove unmodified status.")
	}

	if status.Staged.Contains(untrackedPath, Untracked) ||
		!status.Unstaged.Contains(untrackedPath, Untracked) {
		t.Errorf("Failed to remove unmodified status.")
	}

	if !status.Staged.Contains(onlyStagedPath, Modified) ||
		status.Unstaged.Contains(onlyStagedPath, Modified) {
		t.Errorf("File should be included in staged but not in unstaged.")
	}

	if status.Staged.Contains(onlyUnstagedPath, Modified) ||
		!status.Unstaged.Contains(onlyUnstagedPath, Modified) {
		t.Errorf("File should be included in unstaged but not in staged.")
	}

	if !status.Staged.Contains(bothPath, Modified) ||
		!status.Unstaged.Contains(bothPath, Modified) {
		t.Errorf("File should be included in both staged and unstaged.")
	}
}
