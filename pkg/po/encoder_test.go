package po

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoderEmptyFile(t *testing.T) {
	var buff bytes.Buffer
	enc := NewEncoder(&buff)
	err := enc.Encode(&File{})
	require.NoError(t, err)
	assert.Empty(t, buff.Bytes())
}

func TestEncoder_Header(t *testing.T) {
	var buff bytes.Buffer
	enc := NewEncoder(&buff)
	err := enc.Encode(&File{
		Header: &Header{
			Comment: &Comment{
				Translator: "My Comment",
			},
			ProjectIDVersion:        "poedit",
			ReportMsgidBugsTo:       "help@poedit.net",
			POTCreationDate:         "2021-06-03 11:36+0200",
			PORevisionDate:          "2022-05-01 15:22+0200",
			LastTranslator:          "Florian Vogt",
			LanguageTeam:            "English, United Kingdom",
			Language:                "en_GB",
			MimeVersion:             "1.0",
			ContentType:             "text/plain; charset=UTF-8",
			ContentTransferEncoding: "8bit",
			PluralForms:             "nplurals=2; plural=(n != 1);",
			XGenerator:              "Poedit 3.0.1",
			UnknownFields: map[string]string{
				"X-Crowdin-File": "/locales/poedit.pot",
			},
		},
	})
	require.NoError(t, err)

	want := `# My Comment
msgid ""
msgstr ""
"Project-Id-Version: poedit\n"
"Report-Msgid-Bugs-To: help@poedit.net\n"
"POT-Creation-Date: 2021-06-03 11:36+0200\n"
"PO-Revision-Date: 2022-05-01 15:22+0200\n"
"Last-Translator: Florian Vogt\n"
"Language-Team: English, United Kingdom\n"
"Language: en_GB\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"
"X-Generator: Poedit 3.0.1\n"
"X-Crowdin-File: /locales/poedit.pot\n"

`
	assert.Equal(t, want, buff.String())
}

func TestEncoder_Message(t *testing.T) {

	file := &File{}
	msg := &Message{
		Comment: &Comment{
			Translator: "Test comment",
			Flags:      []string{"fuzzy"},
		},
		Context:  "ctx",
		ID:       "test",
		IDPlural: "test_plural",
		Str:      map[int]string{},
	}
	file.AddMessage(msg)
	var buff bytes.Buffer
	enc := NewEncoder(&buff)

	t.Run("Placeholders for translations are added", func(t *testing.T) {

		want := `# Test comment
#, fuzzy
msgctxt "ctx"
msgid "test"
msgid_plural "test_plural"
msgstr[0] ""
msgstr[1] ""

`
		err := enc.Encode(file)
		require.NoError(t, err)
		assert.Equal(t, want, buff.String())
	})

	// Add reference
	msg.AddReference(&Reference{Path: "z"})

	// Set translations
	msg.Str[0] = "a"
	msg.Str[1] = "b"
	msg.Str[2] = "c"

	t.Run("Translations are given out", func(t *testing.T) {
		want := `# Test comment
#: z
#, fuzzy
msgctxt "ctx"
msgid "test"
msgid_plural "test_plural"
msgstr[0] "a"
msgstr[1] "b"
msgstr[2] "c"

`
		buff.Reset()

		err := enc.Encode(file)
		require.NoError(t, err)
		assert.Equal(t, want, buff.String())
	})

	// Add messsage
	o := NewMessage()
	o.ID = "mytest"
	o.Comment.Extracted = "TRANSLATORS: This is a unit test"
	o.AddReference(&Reference{Path: "a", Line: 10})
	file.AddMessage(o)

	t.Run("Messages are sorted", func(t *testing.T) {

		want := `#. TRANSLATORS: This is a unit test
#: a:10
msgid "mytest"
msgstr ""

# Test comment
#: z
#, fuzzy
msgctxt "ctx"
msgid "test"
msgid_plural "test_plural"
msgstr[0] "a"
msgstr[1] "b"
msgstr[2] "c"

`

		buff.Reset()
		err := enc.Encode(file)
		require.NoError(t, err)
		assert.Equal(t, want, buff.String())
	})
}

func TestEncoderDisableWriteReference(t *testing.T) {
	t.Run("References are written as standard", func(t *testing.T) {
		file := &File{
			Header: &Header{Comment: &Comment{Translator: "Copyright"}},
		}
		msg := &Message{
			ID:       "test",
			IDPlural: "test_plural",
		}
		msg.AddReference(&Reference{Path: "z"})
		msg.AddReference(&Reference{Path: "x", Line: 5})
		file.AddMessage(msg)
		var buff bytes.Buffer
		enc := NewEncoder(&buff)
		enc.SetWriteEmptyHeader(false)
		enc.SetWriteReferences(true)

		err := enc.Encode(file)
		require.NoError(t, err)

		want := `# Copyright

#: x:5 z
msgid "test"
msgid_plural "test_plural"
msgstr[0] ""
msgstr[1] ""

`
		assert.Equal(t, want, buff.String())
	})

	t.Run("Disable do not write the reference", func(t *testing.T) {
		file := &File{
			Header: &Header{Comment: &Comment{Translator: "Copyright"}},
		}
		msg := &Message{
			ID:       "test",
			IDPlural: "test_plural",
		}
		msg.AddReference(&Reference{Path: "z"})
		file.AddMessage(msg)
		var buff bytes.Buffer
		enc := NewEncoder(&buff)
		enc.SetWriteHeader(false)
		enc.SetWriteReferences(false)

		err := enc.Encode(file)
		require.NoError(t, err)

		want := `msgid "test"
msgid_plural "test_plural"
msgstr[0] ""
msgstr[1] ""

`
		assert.Equal(t, want, buff.String())
	})
}
