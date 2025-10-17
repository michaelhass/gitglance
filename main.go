package main

import (
	"fmt"
	"os"

	"github.com/michaelhass/gitglance/internal/app"
)

func main() {
	var opts []app.Option

	for _, arg := range os.Args {
		if arg == "debug" {
			opts = append(opts, app.WithDebugLogger())
			break
		}
	}

	if err := app.Launch(opts...); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(0)
	}
}
