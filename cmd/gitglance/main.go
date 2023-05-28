package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/ui"
)

func main() {
	err := ui.LaunchApp(
		ui.LaunchOptions{},
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
