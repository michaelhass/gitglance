package git

import (
	"testing"
)

func TestFileStatusFromLineSimple(t *testing.T) {
	var (
		line     = "M some/path"
		expectFS = FileStatus{
			Code: Modified,
			Path: "some/path",
		}
		gotFS FileStatus
		err   error
	)

	gotFS, err = fileStatusFromLine(line)
	if err != nil {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}

	if gotFS.Code != expectFS.Code {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}

	if gotFS.Path != expectFS.Path {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}
}

func TestFileStatusFromLineTrimmingWhitespaces(t *testing.T) {
	var (
		line     = "   M some/path   "
		expectFS = FileStatus{
			Code: Modified,
			Path: "some/path",
		}
		gotFS FileStatus
		err   error
	)

	gotFS, err = fileStatusFromLine(line)
	if err != nil {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}

	if gotFS.Code != expectFS.Code {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}

	if gotFS.Path != expectFS.Path {
		t.Error("Unable to retrieve FileStatus from line: ", line)
	}
}

func TestFileStatusFromLineMissingPath(t *testing.T) {
	var (
		line = "M"
		err  error
	)

	_, err = fileStatusFromLine(line)
	if err == nil {
		t.Error("Missing path. Expected error for line: ", line)
	}
}

func TestFileStatusFromLineUnknownStatus(t *testing.T) {
	var (
		line = "Y some/path"
		err  error
	)

	_, err = fileStatusFromLine(line)
	if err == nil {
		t.Error("Unknown status code. Expected error for line: ", line)
	}
}

func TestFileStatusFromLineEmpty(t *testing.T) {
	var (
		line = ""
		err  error
	)

	_, err = fileStatusFromLine(line)
	if err == nil {
		t.Error("Empty line. Expected error")
	}
}
