package git

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

type Repository interface {
	Worktree() (Worktree, error)
}

type RepositoryOpt struct {
	ImplType ImplType
}

type ImplType byte

const (
	GoGit ImplType = iota
)

var (
	errNotImplemented = errors.New("not implemented")
)

func OpenRepository(path string, opt RepositoryOpt) (Repository, error) {
	switch opt.ImplType {
	case GoGit:
		r, err := git.PlainOpen(path)
		if err != nil {
			return nil, err
		}
		return &goGitRepository{repo: r}, nil
	default:
		return nil, errNotImplemented
	}
}

type goGitRepository struct {
	repo *git.Repository
}

func (r *goGitRepository) Worktree() (Worktree, error) {
	wt, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}
	return newGoGitWorkTree(wt), nil
}
