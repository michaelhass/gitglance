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

	sp, err := git.NewStatusProvider(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	status, err := sp.CurrentStatus()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Printf("%v", status)
}
