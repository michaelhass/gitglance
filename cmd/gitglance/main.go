package main

import (
	"os"

	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui"
)

func main() {
	args := os.Args[1:]

	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "."
	}

	ui.LaunchApp(
		ui.LaunchOptions{
			Path:    path,
			RepoOpt: git.RepositoryOpt{ImplType: git.GoGit},
		},
	)
}
