package editor

import (
	"os"
	"os/exec"
)

func OpenFileCmd(path string) *exec.Cmd {
	command := "vi"
	if visual := os.Getenv("VISUAL"); len(visual) > 0 {
		command = visual
	} else if editor := os.Getenv("EDITOR"); len(visual) > 0 {
		command = editor
	}
	return exec.Command(command, path)
}
