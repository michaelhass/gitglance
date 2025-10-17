package git

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	stashEntryIdxRegexPattern    = `stash@{([0-9]+)}`
	stashEntryComponentSeparator = ": "
)

var (
	missingStashIdxErr                  = errors.New("No stash index found")
	missingStashMsgErr                  = errors.New("No stash message found")
	defaultStashIdxRegex *regexp.Regexp = regexp.MustCompile(stashEntryIdxRegexPattern)
)

type stashBuilder struct {
	stashEntryIdxRegex *regexp.Regexp
}

func newDefaultStashBuilder() stashBuilder {
	return stashBuilder{
		stashEntryIdxRegex: defaultStashIdxRegex,
	}
}

func (b stashBuilder) makeStashFromMultilineText(text string) (Stash, error) {
	lines := strings.Split(text, "\n")
	return b.makeStashFromLines(lines...)
}

func (b stashBuilder) makeStashFromLines(lines ...string) (Stash, error) {
	var (
		entries []StashEntry
		err     error
	)

	for _, line := range lines {
		idx, idxErr := b.getStashEntryIdxFromLine(line)
		if idxErr != nil {
			err = idxErr
			break
		}
		msg, msgErr := b.getStashEntryMsgFromLine(line)
		if msgErr != nil {
			err = msgErr
			break
		}

		entries = append(entries, StashEntry{idx: idx, msg: msg})
	}

	return entries, err
}

func (b stashBuilder) getStashEntryIdxFromLine(line string) (int, error) {
	matches := b.stashEntryIdxRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return -1, missingStashIdxErr
	}
	idx, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, missingStashIdxErr
	}
	return idx, nil
}

func (b stashBuilder) getStashEntryMsgFromLine(line string) (string, error) {
	idx := strings.Index(line, stashEntryComponentSeparator)
	if idx == -1 {
		return "", missingStashMsgErr
	}
	msgStartIdx := idx + len(stashEntryComponentSeparator)
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

type Stash []StashEntry
