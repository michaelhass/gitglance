package git

import (
	"fmt"
	"strings"
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

type StatusCode byte

const (
	Unmodified         StatusCode = ' '
	Untracked          StatusCode = '?'
	Modified           StatusCode = 'M'
	Added              StatusCode = 'A'
	Deleted            StatusCode = 'D'
	Renamed            StatusCode = 'R'
	Copied             StatusCode = 'C'
	UpdatedButUnmerged StatusCode = 'U'
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
