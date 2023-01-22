package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Status struct {
	Unstaged []FileStatus
	Staged   []FileStatus
}

func (s Status) contains(src []FileStatus, path string, code StatusCode) bool {
	for _, fs := range src {
		if fs.Path == path && fs.Code == code {
			return true
		}
	}
	return false
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
