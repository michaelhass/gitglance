package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/ui/app"
)

func main() {
	opts := app.LaunchOptions{IsDebug: false}
	if err := app.Launch(opts); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(0)
	}
}
