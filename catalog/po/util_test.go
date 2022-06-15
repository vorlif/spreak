package po

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_Newline(t *testing.T) {
	raw := `"test\n"`
	decoded := DecodePoString(raw)
	assert.Equal(t, "test\n", decoded)
}

func TestDecodeMultilineString(t *testing.T) {
	raw := `""
"this is a \n"
"multiline string"`
	decoded := DecodePoString(raw)
	assert.Equal(t, "this is a \nmultiline string", decoded)

	raw = `"Plural-Forms: nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 "
"&& n%100<=10 ? 3 : n%100>=11 && n%100<=99 ? 4 : 5);\n"`
	decoded = DecodePoString(raw)

	assert.Equal(t, "Plural-Forms: nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 && n%100<=10 ? 3 : n%100>=11 && n%100<=99 ? 4 : 5);\n", decoded)
}

func TestEncodeMultilineString(t *testing.T) {
	t.Run("default cases", func(t *testing.T) {
		tests := []struct {
			name string
			raw  string
			want string
		}{
			{
				name: "Long lines will be wrapped",
				raw: `
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.
At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, 
no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr,      sed diam nonumy eirmod tempor
`,
				want: `""
"\n"
"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, "
"sed diam nonumy eirmod tempor invidunt ut labore et dolore "
"magna aliquyam erat, sed diam voluptua.\n"
"At vero eos et accusam et justo duo dolores et ea rebum. "
"Stet clita kasd gubergren, \n"
"no sea takimata sanctus est Lorem ipsum dolor sit amet. \n"
"Lorem ipsum dolor sit amet, consetetur sadipscing "
"elitr,      sed diam nonumy eirmod tempor\n"`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				get := EncodePoString(tt.raw, 60)
				assert.Equal(t, tt.want, get)
			})
		}
	})
}

func TestDecodePoString(t *testing.T) {
	t.Run("header", func(t *testing.T) {
		assert.Equal(t, poStrDecode, DecodePoString(poStrEncode))
	})

	for i, tt := range decodeEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := DecodePoString(tt.poLines)
			assert.Equal(t, tt.code, actual, "id=%d", i)
		})
	}
}

func TestEncodePoString(t *testing.T) {
	t.Run("header", func(t *testing.T) {
		assert.Equal(t, poStrEncodeStd, EncodePoString(poStrDecode, 60))
	})

	for i, tt := range decodeEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EncodePoString(tt.code, 60)
			assert.Equal(t, tt.poLines, actual, "id=%d wrap=enabled", i)
			actual = EncodePoString(tt.code, -1)
			assert.Equal(t, tt.poLines, actual, "id=%d wrap=disabled", i)
		})
	}
}

func TestEncodePoStringWrap(t *testing.T) {

	cases := []struct {
		Input, Output string
		Lim           int
	}{
		// A simple word passes through.
		{
			"foo",
			`"foo"`,
			4,
		},
		// A single word that is too long passes through.
		// We do not break words.
		{
			"foobarbaz",
			`"foobarbaz"`,
			4,
		},
		// Lines are broken at whitespace.
		{
			"foo bar baz",
			`"foo "
"bar "
"baz"`,
			4,
		},
		// Lines are broken at whitespace, even if words
		// are too long. We do not break words.
		{
			"foo bars bazzes",
			`"foo "
"bars "
"bazzes"`,
			4,
		},
		// A word that would run beyond the width is wrapped.
		{
			"fo sop",
			`"fo "
"sop"`,
			4,
		},
		// An explicit line break at the end of the input is preserved.
		{
			"foo bar baz\n",
			`"foo "
"bar "
"baz\n"`,
			4,
		},
		// Multi-byte characters
		{
			"\u2584 \u2584 \u2584 \u2584",
			"\"\u2584 \u2584 \"\n\"\u2584 \u2584\"",
			4,
		},
	}

	for _, tc := range cases {
		actual := strings.Join(encodePoStringWithWrap(tc.Input, tc.Lim), "\n")
		assert.Equal(t, tc.Output, actual)
	}

}

const poStrEncode = `# noise
123456789
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
"X-Crowdin-Project: poedit\n"
"X-Crowdin-Project-ID: 53425\n"
"X-Crowdin-Language: en-GB\n"
"X-Crowdin-File: /locales/poedit.pot\n"
"X-Crowdin-File-ID: 3\n"
"X-Generator: Poedit 3."
"0.1\n"
>>
123456???
`

const poStrDecode = `Project-Id-Version: poedit
Report-Msgid-Bugs-To: help@poedit.net
POT-Creation-Date: 2021-06-03 11:36+0200
PO-Revision-Date: 2022-05-01 15:22+0200
Last-Translator: Florian Vogt
Language-Team: English, United Kingdom
Language: en_GB
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit
Plural-Forms: nplurals=2; plural=(n != 1);
X-Crowdin-Project: poedit
X-Crowdin-Project-ID: 53425
X-Crowdin-Language: en-GB
X-Crowdin-File: /locales/poedit.pot
X-Crowdin-File-ID: 3
X-Generator: Poedit 3.0.1
`

const poStrEncodeStd = `""
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
"X-Crowdin-Project: poedit\n"
"X-Crowdin-Project-ID: 53425\n"
"X-Crowdin-Language: en-GB\n"
"X-Crowdin-File: /locales/poedit.pot\n"
"X-Crowdin-File-ID: 3\n"
"X-Generator: Poedit 3.0.1\n"`

var decodeEncodeTests = []struct {
	name    string
	code    string
	poLines string
}{
	{
		name:    "Quotation marks are added",
		code:    `Extracted string`,
		poLines: `"Extracted string"`,
	},
	{
		name:    "Intentional multiple quotation marks are preserved",
		code:    "\"Extracted string\"",
		poLines: "\"\\\"Extracted string\\\"\"",
	},
	{
		name:    "Intentional backquotes are preserved",
		code:    "`Extracted string`",
		poLines: "\"`Extracted string`\"",
	},
	{
		name: "Multiline text are formatted correctly",
		code: "This is an multiline\nstring",
		poLines: `""
"This is an multiline\n"
"string"`,
	},
	{
		name: "backquoted newline is converted to newline",
		code: `This is an multiline
string`,
		poLines: `""
"This is an multiline\n"
"string"`,
	},
	{
		name:    "Single line with a newline at the end remains a single line",
		code:    "Single line with newline\n",
		poLines: `"Single line with newline\n"`,
	},
	{
		name: "Last newline does not start a new line",
		code: "Multiline\nwith\nnewlines\n",
		poLines: `""
"Multiline\n"
"with\n"
"newlines\n"`,
	},
	{
		name:    "Empty string formatted",
		code:    "",
		poLines: "\"\"",
	},
	{
		name: "special chars",
		code: `{}()[]special	"chars"	\`,
		poLines: "\"{}()[]special\\t\\\"chars\\\"\\t\\\\\"",
	},
	{
		name:    "backslash",
		code:    "Test \\\\ \"id\"\t\\",
		poLines: "\"Test \\\\\\\\ \\\"id\\\"\\t\\\\\"",
	},
	{
		name:    "ascii",
		code:    "\a\f\v",
		poLines: "\"\\a\\f\\v\"",
	},
	{
		name:    "utf8",
		code:    "äüö€$翻訳ファイル",
		poLines: `"äüö€$翻訳ファイル"`,
	},
}
