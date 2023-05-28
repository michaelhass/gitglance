package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type statusOption struct {
	isPorcelain     bool
	isNULTerminated bool
	hasBranch       bool
}

func newStatusCmd(opt statusOption) *exec.Cmd {
	args := []string{"status"}

	if opt.isPorcelain {
		args = append(args, "--porcelain")
	}

	if opt.isNULTerminated {
		args = append(args, "-z")
	}

	if opt.hasBranch {
		args = append(args, "-b")
	}

	return exec.Command("git", args...)
}

func main() {
	out, err := newStatusCmd(statusOption{
		isPorcelain:     true,
		isNULTerminated: true,
	}).Output()

	if err != nil {
		fmt.Println("could not execute status cmd")
		os.Exit(0)
	}

	statusString := string(out)
	components := strings.Split(statusString, "\000")

	for i := 0; i < len(components); i++ {
		component := components[i]
		if len(component) < 3 {
			continue
		}
		changes := component[0:2]
		path := component[3:]

		fmt.Printf("changes '%v'\n", changes)
		fmt.Printf("path '%v'\n", path)

		if strings.Contains(changes, "R") {
			i += 1
			fmt.Printf("Renamed to '%s'\n", components[i])
		}

		fmt.Println()
	}
}
