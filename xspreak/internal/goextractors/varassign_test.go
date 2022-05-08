package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablesExtractorRun(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewVariablesExtractor())
	assert.NotEmpty(t, issues)

	want := []string{
		"Bob", "Bobby", "application", "john", "doe", "assign function param", "struct attr assign",
		"Newline remains\n", "This is an\nmultiline string",
	}
	got := collectIssueStrings(issues)
	assert.ElementsMatch(t, want, got)
}
