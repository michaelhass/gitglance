package git

import (
	"fmt"
	"strings"
	"testing"
)

func TestWorkTreeStatusBranch(t *testing.T) {
	var (
		branch              = "some_branch"
		branchComponent     = fmt.Sprintf("## %s.../origin/%s", branch, branch)
		stagedFileComponent = "A  cmd/playground/main.go"
		out                 = statusOutputFromComponents([]string{
			branchComponent,
			stagedFileComponent,
		})
		workTreeStatus, err = readWorkTreeStatusFromOutput(out)
	)

	if err != nil {
		t.Error(err)
	}

	if branch != workTreeStatus.CleanedBranchName {
		t.Errorf(
			"Branch not read. Expected '%s' got '%s'",
			branch,
			workTreeStatus.Branch,
		)
	}

	if len(workTreeStatus.StagedFiles()) != 1 {
		t.Errorf(
			"Failed to read staged files. Expected one staged file,  got '%v'",
			workTreeStatus.StagedFiles(),
		)
	}
}

func TestWorkTreeStatusRenamed(t *testing.T) {
	var (
		changes             = "R "
		path                = "some/path/new_name.txt"
		oldPathComponent    = "some/path/old_bame.txt"
		renamedComponent    = fmt.Sprintf("%s %s -> %s", changes, oldPathComponent, path)
		out                 = statusOutputFromComponents([]string{renamedComponent})
		workTreeStatus, err = readWorkTreeStatusFromOutput(out)
	)

	if err != nil {
		t.Error(err)
	}

	files := workTreeStatus.Filter(func(fs FileStatus) bool {
		return fs.StagedStatusCode == Renamed
	})

	if len(files) == 0 {
		t.Error("Did not create renamed FileStatus")
	}

	file := files[0]

	if file.Path != path {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}

	if file.Extra != oldPathComponent {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}
}

func TestStagedFileStatus(t *testing.T) {
	var (
		changes   = "A "
		path      = "some/path/file.txt"
		component = fmt.Sprintf("%s %s", changes, path)
		file, err = readFileStatusFromOutputComponent(component)
	)

	if err != nil {
		t.Error(err)
	}

	if file.HasUnstagedChanges() {
		t.Error("Failed to read status. Expected no unstaged changes")
	}

	if !file.HasStagedChanges() {
		t.Error("Failed to read status, Expected staged changes.")
	}

	if file.StagedStatusCode != Added {
		t.Errorf(
			"Failed to read status. Expexted 'Added' got '%s'",
			string(file.StagedStatusCode),
		)
	}

	if file.Path != path {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}
}

func TestUnstagedFileStatus(t *testing.T) {
	var (
		changes   = " M"
		path      = "some/path/file.txt"
		component = fmt.Sprintf("%s %s", changes, path)
		file, err = readFileStatusFromOutputComponent(component)
	)

	if err != nil {
		t.Error(err)
	}

	if !file.HasUnstagedChanges() {
		t.Error("Failed to read status. Expected no unstaged changes")
	}

	if file.HasStagedChanges() {
		t.Error("Failed to read status, Expected staged changes.")
	}

	if file.UnstagedStatusCode != Modified {
		t.Errorf(
			"Failed to read status. Expexted 'Added' got '%s'",
			string(file.UnstagedStatusCode),
		)
	}

	if file.Path != path {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}
}

func TestUnmodifiedFileStatus(t *testing.T) {
	var (
		changes   = "  "
		path      = "some/path/file.txt"
		component = fmt.Sprintf("%s %s", changes, path)
		file, err = readFileStatusFromOutputComponent(component)
	)

	if err != nil {
		t.Error(err)
	}

	if !file.IsUnmodified() {
		t.Error("Failed to read status. Expected no changes.")
	}

	if file.Path != path {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}
}

func TestUntrackedFileStatus(t *testing.T) {
	var (
		changes   = "? "
		path      = "some/path/file.txt"
		component = fmt.Sprintf("%s %s", changes, path)
		file, err = readFileStatusFromOutputComponent(component)
	)

	if err != nil {
		t.Error(err)
	}

	if !file.IsUntracked() {
		t.Error("Failed to read status. Expected to be untracked.")
	}

	if file.Path != path {
		t.Errorf(
			"Failed to read path. Expexted '%s' got '%s'",
			path,
			file.Path,
		)
	}
}

func statusOutputFromComponents(components []string) string {
	return strings.Join(components, "\n")
}
