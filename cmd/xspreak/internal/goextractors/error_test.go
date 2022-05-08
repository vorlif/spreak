package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorExtractor(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewCommentsExtractor(), NewErrorExtractor())
	assert.NotEmpty(t, issues)

	got := collectIssueStrings(issues)
	want := []string{"global error", "errors", "global alias error",
		"errors", "local error", "errors", "local alias error", "errors", "return error", "errors"}
	assert.ElementsMatch(t, want, got)
}
