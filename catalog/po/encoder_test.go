package po

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleMarshal() {
	f := NewFile()
	f.Header = PlaceholderHeader("Spreak", "Florian Vogt", "info@example.org")
	msg := NewMessage()
	msg.ID = "Hello"
	msg.Str[0] = "Hallo"
	f.AddMessage(msg)

	doc, err := Marshal(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(doc))
}

func ExampleNewEncoder() {
	f := NewFile()
	msg := NewMessage()
	msg.ID = "Hello"
	msg.Str[0] = "Hallo"
	f.AddMessage(msg)

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff)
	enc.SetWriteEmptyHeader(false)
	err := enc.Encode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(buff.String())
	// Output:
	// msgid "Hello"
	// msgstr "Hallo"

}

func ExampleEncoder_Encode() {
	f := NewFile()
	msg := NewMessage()
	msg.ID = "Car"
	msg.IDPlural = "Cars"
	msg.Str[0] = "Auto"
	msg.Str[1] = "Autos"
	f.AddMessage(msg)

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff)
	enc.SetWriteEmptyHeader(false)
	err := enc.Encode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(buff.String())
	// Output:
	// msgid "Car"
	// msgid_plural "Cars"
	// msgstr[0] "Auto"
	// msgstr[1] "Autos"
}

func TestMarshal(t *testing.T) {
	f := NewFile()
	msg := NewMessage()
	msg.ID = "Hello"
	msg.Str[0] = "Hallo"
	f.AddMessage(msg)

	doc, err := Marshal(f)
	if err != nil {
		panic(err)
	}
	want := `
msgid "Hello"
msgstr "Hallo"`
	assert.Contains(t, string(doc), want)
}

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

	t.Run("Multiline are written correctly", func(t *testing.T) {
		f := NewFile()
		f.Header = nil

		msg = NewMessage()
		msg.ID = "Australian English"
		msg.Str[0] = `ئاۋىستىرالىيە ئىنگلزچىسى
 `
		f.AddMessage(msg)
		buff.Reset()
		err := enc.Encode(f)
		require.NoError(t, err)
		want := `msgid "Australian English"
msgstr ""
"ئاۋىستىرالىيە ئىنگلزچىسى\n"
" "
`
		assert.Equal(t, want, buff.String())
	})

	t.Run("flags encoded correctly", func(t *testing.T) {
		f := NewFile()
		f.Header = nil

		msg = NewMessage()
		msg.ID = "Hello"
		msg.Comment = NewComment()
		msg.Comment.Flags = append(msg.Comment.Flags, "fuzzy", "go-format")
		f.AddMessage(msg)
		buff.Reset()
		err := enc.Encode(f)
		require.NoError(t, err)
		want := `#, fuzzy, go-format
msgid "Hello"
msgstr ""
`
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

func TestEncoder_SetWrapWidth(t *testing.T) {
	file := &File{
		Header: &Header{Comment: &Comment{Translator: "This is an very long line with many words without meaning"}},
	}
	var buff bytes.Buffer
	enc := NewEncoder(&buff)
	enc.SetWrapWidth(15)
	enc.SetWriteEmptyHeader(false)
	err := enc.Encode(file)
	assert.NoError(t, err)

	want := `# This is an very
# long line with
# many words
# without meaning

`
	assert.Equal(t, want, buff.String())
}

func TestDefaultSortFunction(t *testing.T) {
	t.Run("sort_by_id", func(t *testing.T) {
		a := NewMessage()
		a.ID = "Alpha"
		b := NewMessage()
		b.ID = "Alpha"

		assert.Zero(t, DefaultSortFunction(a, b))

		b.ID = "Beta"
		assert.Equal(t, -1, DefaultSortFunction(a, b))

		a.ID = "Charlie"
		assert.Equal(t, 1, DefaultSortFunction(a, b))
	})

	t.Run("sort_by_context", func(t *testing.T) {
		a := NewMessage()

		a.ID = "Alpha"
		a.Context = "ctx"
		b := NewMessage()
		b.ID = "Alpha"
		b.Context = "ctx"

		assert.Zero(t, DefaultSortFunction(a, b))
		a.Context = ""

		assert.Equal(t, -1, DefaultSortFunction(a, b))

		a.Context = "actx"
		assert.Equal(t, -1, DefaultSortFunction(a, b))

		a.Context = "dctx"
		assert.Equal(t, 1, DefaultSortFunction(a, b))
	})

	t.Run("sort_by_references", func(t *testing.T) {
		a := NewMessage()
		a.ID = "Alpha"
		aRef := &Reference{Path: "a.go"}
		a.AddReference(aRef)
		b := NewMessage()
		b.ID = "Alpha"
		bRef := &Reference{Path: "a.go"}
		b.AddReference(bRef)

		assert.Zero(t, DefaultSortFunction(a, b))

		bRef.Path = "b.go"
		assert.Equal(t, -1, DefaultSortFunction(a, b))

		aRef.Path = "c.go"
		assert.Equal(t, 1, DefaultSortFunction(a, b))

		aRef.Path = "a.go"
		aRef.Line = 1
		bRef.Path = "a.go"
		bRef.Line = 2

		assert.Equal(t, -1, DefaultSortFunction(a, b))
		aRef.Line = 3
		assert.Equal(t, 1, DefaultSortFunction(a, b))
		bRef.Line = 3

		aRef.Column = 30
		bRef.Column = 20
		assert.Equal(t, 1, DefaultSortFunction(a, b))
		aRef.Column = 10
		assert.Equal(t, -1, DefaultSortFunction(a, b))

		a.AddReference(&Reference{Path: "z.go"})
		b.AddReference(&Reference{Path: "d.go"})

		assert.Equal(t, -1, DefaultSortFunction(a, b))
	})
}
