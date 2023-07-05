package git

func Status() (WorkTreeStatus, error) {
	return loadWorkTreeStatus()
}

func StageFile(path string) error {
	cmd := newGitCommand("add", path)
	return cmd.run()
}

func UnstageFile(path string) error {
	return newGitCommand("restore", "--staged", path).run()
}

func Diff(opt DiffOptions) (string, error) {
	return newDiffCmd(opt).output()
}

func Commit(msg string) error {
	return newGitCommand("commit", "-m", msg).run()
}
