package po

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Simple(t *testing.T) {
	t.Run("test parse empty content", func(t *testing.T) {
		file, err := Parse([]byte{})
		assert.Error(t, err)
		assert.Nil(t, file)
	})

	t.Run("test parse nil", func(t *testing.T) {
		file, err := Parse(nil)
		assert.Error(t, err)
		assert.Nil(t, file)
	})

	t.Run("test parse message", func(t *testing.T) {
		content := `
msgid "id"
msgstr "ID"`
		file, err := Parse([]byte(content))
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.NotNil(t, file.Header)
		assert.Equal(t, len(file.Messages), 1)
		if assert.Contains(t, file.Messages, "") && assert.Contains(t, file.Messages[""], "id") {
			msg := file.Messages[""]["id"]
			assert.Equal(t, "id", msg.ID)
			assert.Equal(t, "ID", msg.Str[0])
		}
	})

	t.Run("test duplicated message id", func(t *testing.T) {
		content := `
msgid "id"
msgstr "ID"

msgid "id"
msgstr ""`
		f, err := Parse([]byte(content))
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test invalid plural index", func(t *testing.T) {
		content := `
msgid "id"
msgid_plural "id_plural"
msgstr[0] "ID"
msgstr[a] "ID"`

		f, err := Parse([]byte(content))
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test multiple headers", func(t *testing.T) {
		content := `
msgid ""
msgstr ""

msgid ""
msgstr ""`
		f, err := Parse([]byte(content))
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test no messages comment", func(t *testing.T) {
		content := "# Comment"
		f, err := Parse([]byte(content))
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}

func TestParse_Header(t *testing.T) {
	testHeader := testHeaderStartComment + testHeaderBody
	file, err := ParseString(testHeader)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.NotNil(t, file.Header)
	assert.Empty(t, file.Messages)

	t.Run("parse comment", func(t *testing.T) {
		header := `# This file is distributed under the same license as the Django package.
#
# Translators:
# F Wolff <friedel@translate.org.za>, 2019-2020
# Stephen Cox <stephencoxmail@gmail.com>, 2011-2012
# unklphil <villiers.strauss@gmail.com>, 2014,2019
msgid ""
msgstr ""
"Project-Id-Version: django\n"
"Report-Msgid-Bugs-To: \n"
"POT-Creation-Date: 2020-05-19 20:23+0200\n"
"PO-Revision-Date: 2020-07-20 19:37+0000\n"
"Last-Translator: F Wolff <friedel@translate.org.za>\n"
"Language-Team: Afrikaans (http://www.transifex.com/django/django/language/"
"af/)\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Language: af\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"`

		file, err := ParseString(header)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.NotNil(t, file.Header)
		assert.Empty(t, file.Messages)

		want := `This file is distributed under the same license as the Django package.

Translators:
F Wolff <friedel@translate.org.za>, 2019-2020
Stephen Cox <stephencoxmail@gmail.com>, 2011-2012
unklphil <villiers.strauss@gmail.com>, 2014,2019`
		assert.Equal(t, want, file.Header.Comment.Translator)

	})
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

func TestParse_PoeditFile(t *testing.T) {
	content, errRead := os.ReadFile("../../testdata/parser/poedit_en_GB.po")
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

func TestParseComment(t *testing.T) {
	t.Run("parse flags", func(t *testing.T) {
		content := `#, python-format
msgid "%(value)s trillion"
msgid_plural "%(value)s trillion"
msgstr[0] "%(value)s Billion"
msgstr[1] "%(value)s Billionen"
`
		file, err := Parse([]byte(content))
		require.NoError(t, err)
		require.NotNil(t, file)
		require.NotNil(t, file.Messages)
		msg := file.GetMessage("", "%(value)s trillion")
		assert.NotNil(t, msg)
		assert.Len(t, msg.Comment.Flags, 1)
		assert.Equal(t, msg.Comment.Flags[0], "python-format")
	})

	t.Run("test parse message prev", func(t *testing.T) {
		content := `
#| msgctxt "prev "
#| "context"
#| msgid "msgid "
#| "prev"
msgid "id"
msgstr "ID"`
		file, err := Parse([]byte(content))
		assert.NoError(t, err)
		require.NotNil(t, file)
		msg := file.GetMessage("", "id")
		require.NotNil(t, msg)
		assert.Equal(t, "msgid prev", msg.Comment.PrevMsgID)
		assert.Equal(t, "prev context", msg.Comment.PrevMsgContext)
	})
}

func TestMustParse(t *testing.T) {
	t.Run("panics on error", func(t *testing.T) {
		content := `this is invalid content for po files`
		f := func() {
			_ = MustParse([]byte(content))
		}
		assert.Panics(t, f)
	})

	t.Run("returns result", func(t *testing.T) {
		content := `
msgid "id"
msgstr "ID"`

		f := MustParse([]byte(content))
		assert.NotNil(t, f)
	})
}
