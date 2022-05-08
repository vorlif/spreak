package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructDefExtractor(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewStructDefExtractor())
	assert.NotEmpty(t, issues)

	got := collectIssueStrings(issues)
	want := []string{
		"global struct msgid", "global struct plural",
		"local struct msgid", "local struct plural",
		"struct msgid arr1", "struct plural arr1",
		"struct msgid arr2", "struct plural arr2",
		"A3", "B3", "C3",
		"A4", "B4", "C4",
	}
	assert.ElementsMatch(t, want, got)
}
