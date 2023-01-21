package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

type StatusProvider struct {
	repo *git.Repository
}

type Status struct {
	Unstaged []FileStatus
	Staged   []FileStatus
}

type FileStatus struct {
	Path string
	Code StatusCode
}

type StatusCode git.StatusCode

const (
	Unmodified         StatusCode = StatusCode(git.Unmodified)
	Untracked          StatusCode = StatusCode(git.Untracked)
	Modified           StatusCode = StatusCode(git.Modified)
	Added              StatusCode = StatusCode(git.Added)
	Deleted            StatusCode = StatusCode(git.Deleted)
	Renamed            StatusCode = StatusCode(git.Renamed)
	Copied             StatusCode = StatusCode(git.Copied)
	UpdatedButUnmerged StatusCode = StatusCode(git.UpdatedButUnmerged)
)

func NewStatusProvider(path string) (*StatusProvider, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &StatusProvider{repo: repo}, nil
}

func (sp *StatusProvider) CurrentStatus() (Status, error) {
	var status Status

	wt, err := sp.repo.Worktree()
	if err != nil {
		return status, err
	}

	ws, err := wt.Status()
	if err != nil {
		return status, err
	}

	for path, fileStatus := range ws {
		if code := StatusCode(fileStatus.Worktree); code != Unmodified {
			status.Unstaged = append(
				status.Unstaged,
				FileStatus{
					Path: path,
					Code: code,
				},
			)
		}
		if code := StatusCode(fileStatus.Staging); code != Unmodified && code != Untracked {
			status.Staged = append(
				status.Staged,
				FileStatus{
					Path: path,
					Code: code,
				},
			)
		}
	}
	return status, nil
}

func (s Status) String() string {
	writeFileStatus := func(sb *strings.Builder, fsList []FileStatus) {
		for _, fs := range fsList {
			sb.WriteString(fmt.Sprintf("[%s] %s\n", string(fs.Code), fs.Path))
		}
	}

	var sb strings.Builder

	sb.WriteString("Unstaged:\n")
	writeFileStatus(&sb, s.Unstaged)
	sb.WriteString("\n")
	sb.WriteString("Staged:\n")
	writeFileStatus(&sb, s.Staged)
	return sb.String()
}
