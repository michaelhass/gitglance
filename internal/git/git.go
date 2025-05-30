// Package git provides easy execution of git commands.
package git

import (
	"fmt"
	"os"
	"strings"
)

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

// StageAll stages all files in the work tree
func StageAll() error {
	return StageFile(".")
}

// UnstageFile unstages a file at the given path.
func UnstageFile(path string) error {
	return newGitCommand("restore", "--staged", path).run()
}

// ResetFile either resets the file to the index or deletes it in case
// it is untracked.
func ResetFile(filePath string, isUntracked bool) error {
	if isUntracked {
		if err := removeFileCmd(filePath).Run(); err != nil && !isExitError(err) {
			return err
		}
		return nil
	}
	return newGitCommand("restore", filePath).run()
}

// UnstageAll unstages all staged files in the work tree.
func UnstageAll() error {
	return UnstageFile(".")
}

// Diff performs a `git diff“ with the given options.
func Diff(opt DiffOptions) (string, error) {
	return newDiffCmd(opt).output()
}

// Commit performs a commit with the given message.
func Commit(msg string) error {
	return newGitCommand("commit", "-m", msg).run()
}

// CurrentBranch returns the name of the current current branch or an error.
func CurrentBranch() (string, error) {
	return newGitCommand("branch", "--show--current").output()
}

// CoreEditorValue returns the currently set local editor for git
// Can be used to direclty open files.
func CoreEditorValue() (string, error) {
	return newGitCommand("config", "core.editor").output()
}

// CoreEditorValue returns the currently set global editor for git
// Can be used to direclty open files.
func CoreGlobalEditorValue() (string, error) {
	return newGitCommand("config", "--global", "core.editor").output()
}

// RootFolder returns the git root folder.
func RootFolder() (string, error) {
	output, err := newGitCommand("rev-parse", "--git-dir").output()
	if err != nil {
		return "", err
	}
	output = strings.ReplaceAll(output, "\n", "")
	return strings.TrimSpace(output), nil
}

// MergeMsg returns the content of the file .git/MERGE_MSG
func MergeMsg() (string, error) {
	folder, err := RootFolder()
	if err != nil {
		return "", err
	}
	mergeMsgPath := fmt.Sprintf("%s/MERGE_MSG", folder)
	mergeFile, err := os.ReadFile(mergeMsgPath)
	if err != nil {
		return "", err
	}
	return string(mergeFile), nil
}

func IsInWorkTree() bool {
	out, err := newGitCommand("rev-parse", "--is-inside-work-tree").output()
	if err != nil {
		return false
	}
	out = strings.ReplaceAll(out, "\n", "")
	return out == "true"
}
