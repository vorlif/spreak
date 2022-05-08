package po

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHeaderStartComment = `# SOME DESCRIPTIVE TITLE.
# Copyright (C) YEAR THE PACKAGE'S COPYRIGHT HOLDER
# This file is distributed under the same license as the Poedit package.
# FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.
#
msgid ""
msgstr ""`

var testHeaderBody = `
"Project-Id-Version: Project 1.0\n"
"Report-Msgid-Bugs-To: help@example.com\n"
"POT-Creation-Date: 2021-06-03 11:36+0200\n"
"PO-Revision-Date: 2022-04-29 23:17+0200\n"
"Last-Translator: Florian Vogt\n"
"Language-Team: \n"
"Language: de\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"
"X-Generator: Poedit 3.0.1\n"
`

func TestHeader(t *testing.T) {

	f, err := ParseString(testHeaderStartComment + testHeaderBody)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NotNil(t, f.Header)
	assert.Empty(t, f.Messages)

	header := f.Header
	a := assert.New(t)
	a.Equal("Project 1.0", header.ProjectIDVersion)
	a.Equal("help@example.com", header.ReportMsgidBugsTo)
	a.Equal("2021-06-03 11:36+0200", header.POTCreationDate)
	a.Equal("2022-04-29 23:17+0200", header.PORevisionDate)
	a.Equal("Florian Vogt", header.LastTranslator)
	a.Equal("", header.LanguageTeam)
	a.Equal("de", header.Language)
	a.Equal("1.0", header.MimeVersion)
	a.Equal("text/plain; charset=UTF-8", header.ContentType)
	a.Equal("8bit", header.ContentTransferEncoding)
	a.Equal("nplurals=2; plural=(n != 1);", header.PluralForms)
	a.Equal("Poedit 3.0.1", header.XGenerator)
}

func TestHeader_MultilineHeader(t *testing.T) {
	newPluralForms := `nplurals=3; plural=(n==1 ? 0 : (n==0 || (n%100>0 && n"
"%100<20)) ? 1 : 2);\n"`
	testHeader := testHeaderStartComment + testHeaderBody
	testHeader = strings.ReplaceAll(testHeader, `nplurals=2; plural=(n != 1);\n"`, newPluralForms)

	f, err := ParseString(testHeader)
	require.NoError(t, err)
	require.NotNil(t, f)
	assert.Empty(t, f.Messages)
	require.NotNil(t, f.Header)

	assert.Equal(t, "nplurals=3; plural=(n==1 ? 0 : (n==0 || (n%100>0 && n%100<20)) ? 1 : 2);", f.Header.PluralForms)

	newPluralForms = `nplurals=4; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n"
"%10<=4 && (n%100<12 || n%100>14) ? 1 : n%10==0 || n%10>=5 && n%10<=9 || n"
"%100>=11 && n%100<=14 ? 2 : 3);\n"`
	testHeader = testHeaderStartComment + testHeaderBody
	testHeader = strings.ReplaceAll(testHeader, `nplurals=2; plural=(n != 1);\n"`, newPluralForms)

	f, err = ParseString(testHeader)
	require.NoError(t, err)
	require.NotNil(t, f)
	assert.Empty(t, f.Messages)
	require.NotNil(t, f.Header)

	assert.Equal(t, "nplurals=4; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<12 || n%100>14) ? 1 : n%10==0 || n%10>=5 && n%10<=9 || n%100>=11 && n%100<=14 ? 2 : 3);", f.Header.PluralForms)
}
