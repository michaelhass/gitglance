package git

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	stashIdxRegexPattern    = `stash@{([0-9]+)}`
	stashComponentSeparator = ": "
)

var (
	missingStashIdxErr                  = errors.New("No stash index found")
	missingStashMsgErr                  = errors.New("No stash message found")
	defaultStashIdxRegex *regexp.Regexp = regexp.MustCompile(stashIdxRegexPattern)
)

type stashEntryBuilder struct {
	stashIdxRegex *regexp.Regexp
}

func newDefaultStashEntryBuilder() stashEntryBuilder {
	return stashEntryBuilder{
		stashIdxRegex: defaultStashIdxRegex,
	}
}

func (b stashEntryBuilder) makeStashEntryFromMultilineText(text string) ([]StashEntry, error) {
	lines := strings.Split(text, "\n")
	return b.makeStashEntryFromLines(lines...)
}

func (b stashEntryBuilder) makeStashEntryFromLines(lines ...string) ([]StashEntry, error) {
	var (
		entries []StashEntry
		err     error
	)

	for _, line := range lines {
		idx, idxErr := b.getStashIdxFromLine(line)
		if idxErr != nil {
			err = idxErr
			break
		}
		msg, msgErr := b.getStashMsgFromLine(line)
		if msgErr != nil {
			err = msgErr
			break
		}

		entries = append(entries, StashEntry{idx: idx, msg: msg})
	}

	return entries, err
}

func (b stashEntryBuilder) getStashIdxFromLine(line string) (int, error) {
	matches := b.stashIdxRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return -1, missingStashIdxErr
	}
	idx, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, missingStashIdxErr
	}
	return idx, nil
}

func (b stashEntryBuilder) getStashMsgFromLine(line string) (string, error) {
	idx := strings.Index(line, stashComponentSeparator)
	if idx == -1 {
		return "", missingStashMsgErr
	}
	msgStartIdx := idx + len(stashComponentSeparator)
	return line[msgStartIdx:], nil
}

type StashEntry struct {
	idx int
	msg string
}

func (s StashEntry) Message() string {
	return s.msg
}

func (s StashEntry) Index() int {
	return s.idx
}
