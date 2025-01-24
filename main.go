package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/ui/app"
)

func main() {
	isDebug := false
	for _, arg := range os.Args {
		isDebug = arg == "debug"
	}
	opts := app.LaunchOptions{IsDebug: isDebug}
	if err := app.Launch(opts); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(0)
	}
}
