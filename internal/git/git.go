package git

import (
	"os/exec"
)

func Status() (WorkTreeStatus, error) {
	var workTreeStatus WorkTreeStatus

	unstagedFiles, err := fileStatusListFromDiff(DiffOption{IsNameStatusOnly: true})
	if err != nil {
		return workTreeStatus, err
	}

	untrackedFilesPaths, err := untrackedFiles()
	if err != nil {
		return workTreeStatus, nil
	}

	untrackedFiles := fileStatusListForUntrackedFiles(untrackedFilesPaths)
	unstagedFiles = append(unstagedFiles, untrackedFiles...)

	stagedFiles, err := fileStatusListFromDiff(DiffOption{IsNameStatusOnly: true, IsStaged: true})
	if err != nil {
		return workTreeStatus, err
	}

	workTreeStatus = NewWorkTreeStatus(unstagedFiles, stagedFiles)
	return workTreeStatus, nil
}

func StageFile(path string) error {
	cmd := exec.Command("git", "add", path)
	return cmd.Run()
}

func UnstageFile(path string) error {
	cmd := exec.Command("git", "restore", "--staged", path)
	return cmd.Run()
}
