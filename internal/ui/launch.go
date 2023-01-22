package ui

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/git"
)

type LaunchOptions struct {
	Path    string
	RepoOpt git.RepositoryOpt
}

func LaunchApp(opt LaunchOptions) {
	fmt.Printf("Open repository at path:'%s'\n", opt.Path)

	repo, err := git.OpenRepository(opt.Path, git.RepositoryOpt{ImplType: git.GoGit})
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	s, err := wt.Status()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Printf("%v", s)
}
