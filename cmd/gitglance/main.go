package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/git"
)

func main() {

	args := os.Args[1:]

	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "."
	}

	fmt.Printf("Open repository at path:'%s'\n", path)
	repo, err := git.OpenRepository(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	status, err := repo.Worktree.Status()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Printf("%v", status)
}
