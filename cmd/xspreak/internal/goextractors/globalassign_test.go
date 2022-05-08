package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalAssignExtractor(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewGlobalAssignExtractor())
	assert.NotEmpty(t, issues)

	got := collectIssueStrings(issues)
	want := []string{"app", "monday", "tuesday", "wednesday", "thursday", "friday"}
	assert.ElementsMatch(t, want, got)
}
