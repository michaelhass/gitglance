package git

import (
	"fmt"
	"strings"
	"unicode"
)

type Status struct {
	Unstaged FileStatusList
	Staged   FileStatusList
}

func newStatus(unstaged, staged FileStatusList) Status {
	return Status{Unstaged: unstaged, Staged: staged}
}

func (l FileStatusList) Contains(path string, code StatusCode) bool {
	for _, fs := range l {
		if fs.Path == path && fs.Code == code {
			return true
		}
	}
	return false
}

type FileStatus struct {
	Path  string
	Extra string // Contains extra information, e.g. previous name
	Code  StatusCode
}

type FileStatusList []FileStatus

type StatusCode byte

const (
	Untracked          StatusCode = '?'
	Modified           StatusCode = 'M'
	Added              StatusCode = 'A'
	Deleted            StatusCode = 'D'
	Renamed            StatusCode = 'R'
	Copied             StatusCode = 'C'
	UpdatedButUnmerged StatusCode = 'U'
)

func (s StatusCode) isValid() bool {
	switch s {
	case Untracked, Modified, Added, Deleted, Renamed, Copied, UpdatedButUnmerged:
		return true
	default:
		return false
	}
}

func fileStatusListForUntrackedFiles(paths []string) FileStatusList {
	var list FileStatusList
	for _, path := range paths {
		list = append(
			list, FileStatus{
				Path: path,
				Code: Untracked,
			},
		)
	}
	return list
}

// Expect format from "git diff --name-status"
func fileStatusListFromDiffString(diff string) (FileStatusList, error) {
	var (
		lines = strings.Split(diff, "\n")
		list  FileStatusList
	)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fs, err := fileStatusFromLine(line)
		if err != nil {
			return nil, err
		}
		list = append(list, fs)
	}
	return list, nil
}

func fileStatusFromLine(line string) (FileStatus, error) {
	var (
		fileStatus  FileStatus
		components  []string
		builder     strings.Builder
		trimmedLine = strings.TrimSpace(line)
	)

	for _, r := range trimmedLine {
		if !unicode.IsSpace(r) {
			builder.WriteRune(r)
			continue
		}

		if builder.Len() == 0 {
			continue
		}

		components = append(components, builder.String())
		builder.Reset()
	}

	if builder.Len() > 0 {
		components = append(components, builder.String())
	}

	if len(components) < 2 {
		return fileStatus,
			fmt.Errorf("unable to read file status from: %s", line)
	}

	statusCode := StatusCode(components[0][0])
	if !statusCode.isValid() {
		return fileStatus,
			fmt.Errorf("invalid statuscode: %s", string(statusCode))
	}

	fileStatus.Code = statusCode
	fileStatus.Path = components[1]

	if len(components) > 2 {
		fileStatus.Extra = components[2]
	}

	return fileStatus, nil
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
