package goextractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncCallExtractor(t *testing.T) {
	issues := runExtraction(t, testdataDir, NewFuncCallExtractor())
	assert.NotEmpty(t, issues)

	want := []string{
		// "f-msgid", "f-plural", "f-context", "f-domain",
		"init", "localizer func call",
		"noop-msgid", "noop-plural", "noop-context", "noop-domain",
		"msgid",
		"msgid-n", "pluralid-n",
		"domain-d", "msgid-d",
		"domain-dn", "msgid-dn", "pluralid-dn",
		"context-pg", "msgid-pg",
		"context-np", "msgid-np", "pluralid-np",
		"domain-dp", "context-dp", "singular-dp",
		"domain-dnp", "context-dnp", "msgid-dnp", "pluralid-dnp",
		"submsgid", "subplural", "foo test",
	}
	got := collectIssueStrings(issues)
	assert.ElementsMatch(t, want, got)
}
