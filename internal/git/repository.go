package git

import (
	"github.com/go-git/go-git/v5"
)

type Repository struct {
	repo     *git.Repository
	worktree *Worktree
}

func OpenRepository(path string) (*Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	wt, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	return &Repository{repo: r, worktree: newWorkTree(wt)}, nil
}

func (r *Repository) Worktree() *Worktree {
	return r.worktree

}
