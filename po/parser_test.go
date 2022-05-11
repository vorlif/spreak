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
	content, errRead := os.ReadFile("../testdata/parser/poedit_en_GB.po")
	require.NoError(t, errRead)

	file, err := Parse(content)
	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotNil(t, file.Header)
	require.NotNil(t, file.Messages)

	t.Run("header", func(t *testing.T) {
		assert.Equal(t, "poedit", file.Header.ProjectIDVersion)
		assert.Equal(t, "nplurals=2; plural=(n != 1);", file.Header.PluralForms)
		assert.Equal(t, "en_GB", file.Header.Language)
	})

	t.Run("context", func(t *testing.T) {
		ctx := "column/row header"
		msgid := "Needs Work"
		if assert.Contains(t, file.Messages, ctx) && assert.Contains(t, file.Messages[ctx], msgid) {
			msg := file.Messages[ctx][msgid]
			assert.Len(t, msg.Str, 1)
			assert.Empty(t, msg.IDPlural)
			assert.Equal(t, "Needs Work", msg.Str[0])
		}
	})

	t.Run("special characters", func(t *testing.T) {
		ctx := ""
		msgid := "Do you want to delete project “%s”?"
		if assert.Contains(t, file.Messages, ctx) && assert.Contains(t, file.Messages[ctx], msgid) {
			msg := file.Messages[ctx][msgid]
			assert.Len(t, msg.Str, 1)
			assert.Empty(t, msg.IDPlural)
			assert.Equal(t, "Do you want to delete project “%s”?", msg.Str[0])
		}
	})

	t.Run("multiline", func(t *testing.T) {
		ctx := ""
		msgid := "Supports all programming languages recognized by GNU gettext tools (PHP, C/C++, C#, Perl, Python, Java, JavaScript and others)."
		translation := "Supports all programming languages recognised by GNU gettext tools (PHP, C/C++, C#, Perl, Python, Java, JavaScript and others)."
		if assert.Contains(t, file.Messages, ctx) && assert.Contains(t, file.Messages[ctx], msgid) {
			msg := file.Messages[ctx][msgid]
			assert.Len(t, msg.Str, 1)
			assert.Empty(t, msg.IDPlural)
			assert.Equal(t, translation, msg.Str[0])
		}
	})

	t.Run("plural", func(t *testing.T) {
		ctx := ""
		msgid := "Pre-translated %u string"
		if assert.Contains(t, file.Messages, ctx) && assert.Contains(t, file.Messages[ctx], msgid) {
			msg := file.Messages[ctx][msgid]
			assert.Len(t, msg.Str, 2)
			assert.Equal(t, "Pre-translated %u strings", msg.IDPlural)
			assert.Equal(t, "Pre-translated %u string", msg.Str[0])
			assert.Equal(t, "Pre-translated %u strings", msg.Str[1])
		}
	})

	t.Run("total count", func(t *testing.T) {
		var count int
		for ctx := range file.Messages {
			count += len(file.Messages[ctx])
		}

		assert.Equal(t, 650, count)
	})

}
