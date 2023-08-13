package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/ui/app"
)

func main() {
	opts := app.LaunchOptions{}
	if err := app.Launch(opts); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
