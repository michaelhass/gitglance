package git

type Repository struct {
	worktree *Worktree
}

func OpenRepository(path string) (*Repository, error) {
	return &Repository{worktree: newWorkTree()}, nil
}

func (r *Repository) Worktree() *Worktree {
	return r.worktree
}
