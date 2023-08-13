package git

import (
	"fmt"
	"strings"
)

// WorkTreeStatus represents the current status of the git work tree.
type WorkTreeStatus struct {
	// The branch at the time of creation.
	// Can be empty if not requested.
	Branch string
	// List of staged and unstaged files
	FileStatusList
}

func loadWorkTreeStatus() (WorkTreeStatus, error) {
	out, err := newStatusCmd(statusOptions{
		isPorcelain:     true,
		isNULTerminated: true,
		hasBranch:       true,
	}).output()

	if err != nil {
		return WorkTreeStatus{}, err
	}

	return readWorkTreeStatusFromOutput(out)
}

func readWorkTreeStatusFromOutput(statusString string) (WorkTreeStatus, error) {
	var (
		components = strings.Split(statusString, nulSeparator)
		branch     string
		files      FileStatusList
		startIdx   = 0
	)

	if firstComponent := components[0]; strings.HasPrefix(firstComponent, branchComponentPrefix) {
		if len(firstComponent) >= 3 {
			branch = firstComponent[3:]
		}
		startIdx = 1
	}

	for i := startIdx; i < len(components); i++ {
		component := components[i]

		file, err := readFileStatusFromOutputComponent(component)

		if err != nil {
			continue
		}

		if file.IsRenamed() {
			// If renamed, next component will be the origianl file name
			i += 1
			file.Extra = components[i]
		}

		files = append(files, file)
	}

	return WorkTreeStatus{Branch: branch, FileStatusList: files}, nil
}

// FileStatus represents the git status of a file.
// It shows the status in the working tree and in the index.
type FileStatus struct {
	Path               string     // The Path of the file
	Extra              string     // Contains extra information, e.g. old name
	UnstagedStatusCode StatusCode // Working tree status
	StagedStatusCode   StatusCode // Index status
}

func readFileStatusFromOutputComponent(component string) (FileStatus, error) {
	var fileStatus FileStatus

	if len(component) < 3 {
		return fileStatus,
			statusError{
				Reason: "Can't read FileStatus. Component is too short.",
			}
	}

	fileStatus.StagedStatusCode = StatusCode(component[0])
	fileStatus.UnstagedStatusCode = StatusCode(component[1])
	fileStatus.Path = component[3:]

	return fileStatus, nil
}

func (fs FileStatus) IsUnmodified() bool {
	return fs.UnstagedStatusCode == Unmodified &&
		fs.StagedStatusCode == Unmodified
}

func (fs FileStatus) HasUnstagedChanges() bool {
	return fs.UnstagedStatusCode != Unmodified
}

func (fs FileStatus) HasStagedChanges() bool {
	return fs.StagedStatusCode != Unmodified
}

func (fs FileStatus) IsUntracked() bool {
	return fs.UnstagedStatusCode == Untracked ||
		fs.StagedStatusCode == Untracked
}

func (fs FileStatus) IsRenamed() bool {
	return fs.StagedStatusCode == Renamed || fs.UnstagedStatusCode == Renamed
}

type FileStatusList []FileStatus

// UnstagedFiles returns all files that have unstaged changes.
func (fl FileStatusList) UnstagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasUnstagedChanges()
	})
}

// StagedFiles returns all files that have staged changes.
func (fl FileStatusList) StagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasStagedChanges() && !fs.IsUntracked()
	})
}

// Filter returns a new FileStatusList that only includes elements accepted by
// the passed filter.
func (fl FileStatusList) Filter(isIncluded func(FileStatus) bool) FileStatusList {
	var result FileStatusList
	for _, fs := range fl {
		if isIncluded(fs) {
			result = append(result, fs)
		}
	}
	return result
}

// StatusCode of a file.
type StatusCode byte

const (
	Unmodified         StatusCode = ' '
	Modified           StatusCode = 'M'
	TypeChanged        StatusCode = 'T'
	Added              StatusCode = 'A'
	Deleted            StatusCode = 'D'
	Renamed            StatusCode = 'R'
	Copied             StatusCode = 'C'
	UpdatedButUnmerged StatusCode = 'U'
	Untracked          StatusCode = '?'
	Ignored            StatusCode = '!'
)

type statusError struct {
	Reason string
}

func (e statusError) Error() string {
	return fmt.Sprint("Git status error:", e.Reason)
}

const (
	nulSeparator          string = "\000"
	branchComponentPrefix string = "##"
)
