package git

import "github.com/go-git/go-git/v5"

type Repository struct {
	repo     *git.Repository
	Worktree *Worktree
}

func OpenRepository(path string) (*Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	wtb, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	return &Repository{
		repo:     r,
		Worktree: newWorkTree(wtb),
	}, nil
}
