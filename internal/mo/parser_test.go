package mo

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_File(t *testing.T) {
	content, errRead := os.ReadFile("../../testdata/parser/poedit_en_GB.mo")
	require.NoError(t, errRead)

	file, err := ParseBytes(content)
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

func TestParse_InvalidByteOrder(t *testing.T) {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.BigEndian, uint32(0x950412df))
	require.NoError(t, err)

	p := newParser(bytes.NewReader(buff.Bytes()))
	require.NotNil(t, p)
	err = p.parseByteOrder()
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidMagicNumber, err)

	p = newParser(bytes.NewReader([]byte{}))
	require.NotNil(t, p)
	err = p.parseByteOrder()
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func TestParse_UseEndian(t *testing.T) {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.BigEndian, uint32(magicLittleEndian))
	require.NoError(t, err)

	p := newParser(bytes.NewReader(buff.Bytes()))
	require.NotNil(t, p)
	err = p.parseByteOrder()
	if assert.NoError(t, err) {
		assert.Equal(t, binary.BigEndian, p.bo)
	}

	buff = &bytes.Buffer{}
	err = binary.Write(buff, binary.LittleEndian, uint32(magicLittleEndian))
	require.NoError(t, err)

	p = newParser(bytes.NewReader(buff.Bytes()))
	require.NotNil(t, p)
	err = p.parseByteOrder()
	if assert.NoError(t, err) {
		assert.Equal(t, binary.LittleEndian, p.bo)
	}
}

func TestParse_EdgeCases(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		f, err := ParseBytes(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("empty input", func(t *testing.T) {
		f, err := ParseBytes([]byte{})
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("unexpected end", func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := binary.Write(buff, binary.BigEndian, uint32(magicLittleEndian))
		require.NoError(t, err)
		f, err := ParseBytes([]byte{})
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}

func TestParseReader(t *testing.T) {
	f, errO := os.Open("../../testdata/parser/poedit_en_GB.mo")
	require.NoError(t, errO)
	defer f.Close()

	file, err := ParseReader(f)
	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotNil(t, file.Header)
	require.NotNil(t, file.Messages)
}
