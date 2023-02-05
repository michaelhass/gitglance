package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/ui"
)

func main() {
	// args := os.Args[1:]

	// var path string
	// if len(args) > 0 {
	// 	path = args[0]
	// } else {
	// 	path = "."
	// }

	err := ui.LaunchApp(
		ui.LaunchOptions{},
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
