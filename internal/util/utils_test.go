package util

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
	src := `This is an
multiline string
`
	want := `""
"This is an\n"
"multiline string\n"`

	get := EncodePoString(src, 60)
	assert.Equal(t, want, get)

	{
		tests := []struct {
			name string
			raw  string
			want string
		}{
			{
				name: "Quotation marks are added",
				raw:  `Extracted string`,
				want: `"Extracted string"`,
			},
			{
				name: "Intentional multiple quotation marks are preserved",
				raw:  "\"Extracted string\"",
				want: "\"\"Extracted string\"\"",
			},
			{
				name: "Intentional backquotes are preserved",
				raw:  "`Extracted string`",
				want: "\"`Extracted string`\"",
			},
			{
				name: "Multiline text are formatted correctly",
				raw:  "This is an multiline\nstring",
				want: `""
"This is an multiline\n"
"string"`,
			},
			{
				name: "backquoted newline is converted to newline",
				raw: `This is an multiline
string`,
				want: `""
"This is an multiline\n"
"string"`,
			},
			{
				name: "Single line with a newline at the end remains a single line",
				raw:  "Single line with newline\n",
				want: `"Single line with newline\n"`,
			},
			{
				name: "Last newline does not start a new line",
				raw:  "Multiline\nwith\nnewlines\n",
				want: `""
"Multiline\n"
"with\n"
"newlines\n"`,
			},
			{
				name: "Empty string formatted",
				raw:  "",
				want: "\"\"",
			},
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

func TestDecodeMoString(t *testing.T) {
	if s := DecodePoString(moStrEncode); s != moStrDecode {
		t.Fatalf(`expect = %s got = %s`, moStrDecode, s)
	}
}

func TestDecodePoString(t *testing.T) {
	if s := DecodePoString(poStrEncode); s != poStrDecode {
		t.Fatalf(`expect = %s got = %s`, poStrDecode, s)
	}
}

func TestEncodePoString(t *testing.T) {
	assert.Equal(t, poStrEncodeStd, EncodePoString(poStrDecode, 60))
}

const moStrEncode = `# noise
123456789
"Project-Id-Version: Poedit 1.5\n"
"Report-Msgid-Bugs-To: poedit@googlegroups.com\n"
"POT-Creation-Date: 2012-07-30 10:34+0200\n"
"PO-Revision-Date: 2013-02-24 21:00+0800\n"
"Last-Translator: Christopher Meng <trans@cicku.me>\n"
"Language-Team: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=1; plural=0;\n"
"X-Generator: Poedit 1.5.5\n"
"TestPoString: abc"
"123\n"
>>
123456???
`

const moStrDecode = `Project-Id-Version: Poedit 1.5
Report-Msgid-Bugs-To: poedit@googlegroups.com
POT-Creation-Date: 2012-07-30 10:34+0200
PO-Revision-Date: 2013-02-24 21:00+0800
Last-Translator: Christopher Meng <trans@cicku.me>
Language-Team: 
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit
Plural-Forms: nplurals=1; plural=0;
X-Generator: Poedit 1.5.5
TestPoString: abc123
`

const poStrEncode = `# noise
123456789
"Project-Id-Version: Poedit 1.5\n"
"Report-Msgid-Bugs-To: poedit@googlegroups.com\n"
"POT-Creation-Date: 2012-07-30 10:34+0200\n"
"PO-Revision-Date: 2013-02-24 21:00+0800\n"
"Last-Translator: Christopher Meng <trans@cicku.me>\n"
"Language-Team: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=1; plural=0;\n"
"X-Generator: Poedit 1.5.5\n"
"TestPoString: abc"
"123\n"
>>
123456???
`

const poStrEncodeStd = `""
"Project-Id-Version: Poedit 1.5\n"
"Report-Msgid-Bugs-To: poedit@googlegroups.com\n"
"POT-Creation-Date: 2012-07-30 10:34+0200\n"
"PO-Revision-Date: 2013-02-24 21:00+0800\n"
"Last-Translator: Christopher Meng <trans@cicku.me>\n"
"Language-Team: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=1; plural=0;\n"
"X-Generator: Poedit 1.5.5\n"
"TestPoString: abc123\n"`

const poStrDecode = `Project-Id-Version: Poedit 1.5
Report-Msgid-Bugs-To: poedit@googlegroups.com
POT-Creation-Date: 2012-07-30 10:34+0200
PO-Revision-Date: 2013-02-24 21:00+0800
Last-Translator: Christopher Meng <trans@cicku.me>
Language-Team: 
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit
Plural-Forms: nplurals=1; plural=0;
X-Generator: Poedit 1.5.5
TestPoString: abc123
`
