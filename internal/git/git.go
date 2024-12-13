// Package git provides easy execution of git commands.
package git

// Status retrieves the current `git status` represented
// by WorkTreeStatus object.
func Status() (WorkTreeStatus, error) {
	return loadWorkTreeStatus()
}

// StageFile stages a file at the given path.
func StageFile(path string) error {
	cmd := newGitCommand("add", path)
	return cmd.run()
}

// UnstageFile unstages a file at the given path.
func UnstageFile(path string) error {
	return newGitCommand("restore", "--staged", path).run()
}

// Diff performs a `git diff“ with the given options.
func Diff(opt DiffOptions) (string, error) {
	return newDiffCmd(opt).output()
}

// Commit performs a commit with the given message.
func Commit(msg string) error {
	return newGitCommand("commit", "-m", msg).run()
}

func CurrentBranch() (string, error) {
	return newGitCommand("branch", "--show--current").output()
}
