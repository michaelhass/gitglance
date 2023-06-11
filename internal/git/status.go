package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type statusOption struct {
	isPorcelain     bool
	isNULTerminated bool
	hasBranch       bool
}

func newStatusCmd(opt statusOption) *exec.Cmd {
	args := []string{"status"}

	if opt.isPorcelain {
		args = append(args, "--porcelain")
	}

	if opt.isNULTerminated {
		args = append(args, "-z")
	}

	if opt.hasBranch {
		args = append(args, "-b")
	}

	return exec.Command("git", args...)
}

type WorkTreeStatus struct {
	Branch string
	FileStatusList
}

func readWorkTreeStatusFromOutput(out []byte) (WorkTreeStatus, error) {
	var (
		statusString = string(out)
		components   = strings.Split(statusString, nulSeparator)
		branch       string
		files        FileStatusList
		startIdx     = 0
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

func loadWorkTreeStatus() (WorkTreeStatus, error) {
	out, err := newStatusCmd(statusOption{
		isPorcelain:     true,
		isNULTerminated: true,
		hasBranch:       true,
	}).Output()

	if err != nil {
		return WorkTreeStatus{}, err
	}

	return readWorkTreeStatusFromOutput(out)
}

type FileStatus struct {
	Path               string
	Extra              string // Contains extra information, e.g. old name
	UnstagedStatusCode StatusCode
	StagedStatusCode   StatusCode
}

func readFileStatusFromOutputComponent(component string) (FileStatus, error) {
	var fileStatus FileStatus

	if len(component) < 3 {
		return fileStatus,
			StatusError{
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

func (fl FileStatusList) UnstagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasUnstagedChanges()
	})
}

func (fl FileStatusList) StagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasStagedChanges() && !fs.IsUntracked()
	})
}

func (fl FileStatusList) Filter(isIncluded func(FileStatus) bool) FileStatusList {
	var result FileStatusList
	for _, fs := range fl {
		if isIncluded(fs) {
			result = append(result, fs)
		}
	}
	return result
}

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

type StatusError struct {
	Reason string
}

func (e StatusError) Error() string {
	return fmt.Sprint("Git status error:", e.Reason)
}

const (
	nulSeparator          string = "\000"
	branchComponentPrefix string = "##"
)
