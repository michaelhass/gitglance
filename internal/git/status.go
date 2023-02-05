package git

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

type WorkTreeStatus struct {
	Unstaged FileStatusList
	Staged   FileStatusList
}

func NewWorkTreeStatus(unstaged, staged FileStatusList) WorkTreeStatus {
	return WorkTreeStatus{Unstaged: unstaged, Staged: staged}
}

func (s WorkTreeStatus) String() string {
	writeFile := func(sb *strings.Builder, fsList []FileStatus) {
		for _, fs := range fsList {
			sb.WriteString(fmt.Sprintf("[%s] %s\n", string(fs.Code), fs.Path))
		}
	}

	var sb strings.Builder
	sb.WriteString("Unstaged:\n")
	writeFile(&sb, s.Unstaged)
	sb.WriteString("\n")
	sb.WriteString("Staged:\n")
	writeFile(&sb, s.Staged)
	return sb.String()
}

type FileStatus struct {
	Path  string
	Extra string // Contains extra information, e.g. previous name
	Code  Code
}

type FileStatusList []FileStatus

func (l FileStatusList) Contains(path string, code Code) bool {
	for _, fs := range l {
		if fs.Path == path && fs.Code == code {
			return true
		}
	}
	return false
}

type Code byte

const (
	Untracked          Code = '?'
	Modified           Code = 'M'
	Added              Code = 'A'
	Deleted            Code = 'D'
	Renamed            Code = 'R'
	Copied             Code = 'C'
	UpdatedButUnmerged Code = 'U'
)

func (s Code) IsValid() bool {
	switch s {
	case Untracked, Modified, Added, Deleted, Renamed, Copied, UpdatedButUnmerged:
		return true
	default:
		return false
	}
}

func fileStatusListFromDiff(opt DiffOption) (FileStatusList, error) {
	diff, err := Diff(opt)
	if err != nil {
		return nil, err
	}
	return fileStatusListFromDiffString(diff)
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

	statusCode := Code(components[0][0])
	if !statusCode.IsValid() {
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

func untrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}
		files = append(files, line)
	}
	return files, nil
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
