package git

import (
	"os/exec"
)

func Status() (WorkTreeStatus, error) {
	return workTreeStatus()
}

func StageFile(path string) error {
	cmd := exec.Command("git", "add", path)
	return cmd.Run()
}

func UnstageFile(path string) error {
	cmd := exec.Command("git", "restore", "--staged", path)
	return cmd.Run()
}

func Diff(opt DiffOption) (string, error) {
	cmd := newDiffCmd(opt)
	out, err := cmd.Output()
	return string(out), err
}
