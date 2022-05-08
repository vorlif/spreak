package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceDefExtractor(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewSliceDefExtractor())
	assert.NotEmpty(t, issues)

	got := collectIssueStrings(issues)
	want := []string{
		"one", "two", "three", "four",
		"six", "seven", "eight", "nine",
		"global struct slice singular", "global struct slice plural", "global ctx",
		"local struct slice singular", "local struct slice plural", "local ctx",
		"global struct slice singular 2", "global struct slice plural 2",
		"local struct slice singular 2", "local struct slice plural 2", "local struct slice ctx 2",
		"A1", "B1", "C1",
		"A2", "B2", "C2",
		"struct slice msgid", "struct slice plural",
	}
	assert.ElementsMatch(t, want, got)
}
