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

func loadWorkTreeStatus() (WorkTreeStatus, error) {
	var (
		workTreeStatus WorkTreeStatus
		out, err       = newStatusCmd(statusOption{
			isPorcelain:     true,
			isNULTerminated: true,
			hasBranch:       true,
		}).Output()
	)

	if err != nil {
		return workTreeStatus, err
	}

	return readWorkTreeStatusFromOutput(out)
}

func readWorkTreeStatusFromOutput(out []byte) (WorkTreeStatus, error) {
	var (
		workTreeStatus WorkTreeStatus
		statusString   = string(out)
		components     = strings.Split(statusString, "\000")
		branch         string
		files          FileStatusList
	)

	if len(components) == 0 {
		return workTreeStatus, StatusError{Reason: "Unable to read git status"}
	}

	branch = components[0]

	for i := 1; i < len(components); i++ {
		component := components[i]

		if len(component) < 3 {
			continue
		}

		file := FileStatus{
			StagedStatusCode:   StatusCode(component[0]),
			UnstagedStatusCode: StatusCode(component[1]),
			Path:               component[3:],
		}

		if file.StagedStatusCode == Renamed ||
			file.UnstagedStatusCode == Renamed {
			// If renamed, next component will be the origianl file name
			i += 1
			file.Extra = components[i]
		}

		files = append(files, file)
	}

	return WorkTreeStatus{Branch: branch, FileStatusList: files}, nil
}

type FileStatus struct {
	Path               string
	Extra              string // Contains extra information, e.g. new name
	UnstagedStatusCode StatusCode
	StagedStatusCode   StatusCode
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

type FileStatusList []FileStatus

func (fl FileStatusList) UnstagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasUnstagedChanges()
	})
}

func (fl FileStatusList) StagedFiles() FileStatusList {
	return fl.Filter(func(fs FileStatus) bool {
		return fs.HasStagedChanges()
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
	ignored            StatusCode = '!'
)

type StatusError struct {
	Reason string
}

func (e StatusError) Error() string {
	return fmt.Sprint("Git status error:", e.Reason)
}
