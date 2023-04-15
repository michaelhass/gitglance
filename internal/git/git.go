package git

import (
	"errors"
	"os/exec"
)

func Status() (WorkTreeStatus, error) {
	return workTreeStatus()
}

func StageFile(path string) error {
	return runCommand(exec.Command("git", "add", path))
}

func UnstageFile(path string) error {
	return runCommand(exec.Command("git", "restore", "--staged", path))
}

func Diff(opt DiffOption) (string, error) {
	cmd := newDiffCmd(opt)
	out, err := cmd.Output()
	if err != nil && !isExitError(err) {
		return "", err
	}
	return string(out), nil
}

func runCommand(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil && !isExitError(err) {
		return err
	}
	return nil
}

func isExitError(err error) bool {
	var ee *exec.ExitError
	return errors.As(err, &ee)
}
