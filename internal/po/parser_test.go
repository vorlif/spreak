package po

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Simple(t *testing.T) {
	test := `
msgid "id"
msgstr "ID"`
	file, err := Parse([]byte(test))
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.NotNil(t, file.Header)
	assert.Equal(t, len(file.Messages), 1)
	if assert.Contains(t, file.Messages, "") && assert.Contains(t, file.Messages[""], "id") {
		msg := file.Messages[""]["id"]
		assert.Equal(t, "id", msg.ID)
		assert.Equal(t, "ID", msg.Str[0])
	}
}

func TestParse_Header(t *testing.T) {
	testHeader := testHeaderStartComment + testHeaderBody
	file, err := ParseString(testHeader)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.NotNil(t, file.Header)
	assert.Empty(t, file.Messages)
}

func TestParse_Context(t *testing.T) {
	test := `
msgctxt "context"
msgid "id"
msgstr "ID"`

	file, err := ParseString(test)
	assert.NoError(t, err)
	if assert.NotNil(t, file) {
		assert.NotNil(t, file.Header)
	}
	if assert.Contains(t, file.Messages, "context") && assert.Contains(t, file.Messages["context"], "id") {
		msg := file.Messages["context"]["id"]
		assert.Equal(t, "context", msg.Context)
	}
}

func TestParse_File(t *testing.T) {
	content, errRead := os.ReadFile("../../testdata/parser/poedit_en_GB.po")
	require.NoError(t, errRead)

	file, err := Parse(content)
	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotNil(t, file.Header)
	require.NotNil(t, file.Messages)

	assert.Equal(t, "poedit", file.Header.ProjectIDVersion)

}
