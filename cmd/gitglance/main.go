package main

import (
	"fmt"
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

	err := ui.LaunchApp(
		ui.LaunchOptions{
			Path:    path,
			RepoOpt: git.RepositoryOpt{ImplType: git.GoGit},
		},
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
