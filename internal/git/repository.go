package git

import (
	"github.com/go-git/go-git/v5"
)

type Repository struct {
	repo *git.Repository
}

func OpenRepository(path string) (*Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repository{repo: r}, nil
}

func (r *Repository) Worktree() (*Worktree, error) {
	wt, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}
	return newWorkTree(wt), nil
}
