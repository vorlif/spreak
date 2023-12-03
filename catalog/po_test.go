package catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

var poFileDomainA = filepath.FromSlash("../testdata/translation-test/de/LC_MESSAGES/a.po")

func getPoCatalogForDomain(t *testing.T, domain string) Catalog {
	data, err := os.ReadFile(poFileDomainA)
	require.NoError(t, err)
	require.NotNil(t, data)

	cat, errC := NewPoDecoder().Decode(language.German, domain, data)
	require.NoError(t, errC)
	require.NotNil(t, cat)
	return cat
}

func TestCatalog_SimplePublicFunctions(t *testing.T) {
	cat := getPoCatalogForDomain(t, "a")
	assert.Equal(t, language.German, cat.Language())

	translation, err := cat.GetTranslation("", "id")
	if assert.NoError(t, err) {
		assert.Equal(t, "ID", translation)
	}

	translation, err = cat.GetPluralTranslation("", "%d day", 1)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Tag", translation)
	}

	translation, err = cat.GetPluralTranslation("", "%d car", 10)
	if assert.Error(t, err) {
		assert.Equal(t, "%d car", translation)
	}

	translation, err = cat.GetTranslation("context", "Test with context")
	if assert.NoError(t, err) {
		assert.Equal(t, "Test mit Context", translation)
	}

	translation, err = cat.GetPluralTranslation("other", "Test with context", 5)
	if assert.Error(t, err) {
		assert.Equal(t, "Test with context", translation)
	}

	translation, err = cat.GetPluralTranslation("context", "%d result", 1)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Ergebniss", translation)
	}

	translation, err = cat.GetPluralTranslation("context", "%d result", 10)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Ergebnisse", translation)
	}

	translation, err = cat.GetPluralTranslation("other", "%d result", 1)
	if assert.Error(t, err) {
		assert.Equal(t, "%d result", translation)
	}

	translation, err = cat.GetPluralTranslation("other", "%d result", 10)
	if assert.Error(t, err) {
		assert.Equal(t, "%d result", translation)
	}
}

const decodePoInvalidPluralForm = `
msgid ""
msgstr ""
"Plural-Forms: invalid-pluralform\n"
`

const decodeTestData = `
msgid "id"
msgstr "ID"

msgid "empty translation"
msgstr ""

msgid "%d day"
msgid_plural "%d days"
msgstr[0] "%d Tag"
msgstr[1] "%d Tage"

msgid "Special\t	\"chars\""
msgstr "Sonder\t	\"zeichen\""
`

func TestPoDecoder(t *testing.T) {
	domain := "my-domain"
	lang := language.German

	dec := NewPoDecoder()
	catl, err := dec.Decode(lang, domain, []byte{})
	assert.Error(t, err)
	assert.Nil(t, catl)

	catl, err = dec.Decode(lang, domain, []byte(decodeTestData))
	assert.NoError(t, err)
	assert.NotNil(t, catl)
	assert.Equal(t, lang, catl.Language())

	translation, err := catl.GetTranslation("", "id")
	assert.NoError(t, err)
	assert.Equal(t, "ID", translation)

	translation, err = catl.GetTranslation("", "xyz")
	assert.Error(t, err)
	assert.Equal(t, "xyz", translation)

	catl, err = dec.Decode(lang, domain, []byte(decodePoInvalidPluralForm+decodeTestData))
	assert.Error(t, err)
	assert.Nil(t, catl)

	t.Run("test special chars", func(t *testing.T) {
		catl, err = dec.Decode(lang, domain, []byte(decodeTestData))
		assert.NoError(t, err)
		require.NotNil(t, catl)

		tr, err := catl.GetTranslation("", "Special		\"chars\"")
		assert.NoError(t, err)
		assert.Equal(t, "Sonder\t\t\"zeichen\"", tr)
	})
}

func TestMoDecoder(t *testing.T) {
	dec := NewMoDecoder()
	decode, err := dec.Decode(language.German, "", []byte{})
	assert.Error(t, err)
	assert.Nil(t, decode)
}

func TestNewPoCLDRDecoder(t *testing.T) {
	dec := NewPoCLDRDecoder()

	catl, err := dec.Decode(language.German, "", []byte(decodeTestData))
	assert.NoError(t, err)
	require.NotNil(t, catl)

	translation, errT := catl.GetPluralTranslation("", "%d day", "1.2")
	assert.NoError(t, errT)
	assert.Equal(t, "%d Tage", translation)
}

func TestCLDRHeader(t *testing.T) {
	header := `
msgid ""
msgstr ""
"Plural-Forms: invalid-pluralform\n"
"X-spreak-use-CLDR: true\n"
`

	dec := NewPoDecoder()
	catl, err := dec.Decode(language.German, "", []byte(header+decodeTestData))
	assert.NoError(t, err)
	require.NotNil(t, catl)
	translation, errT := catl.GetPluralTranslation("", "%d day", "1.2")
	assert.NoError(t, errT)
	assert.Equal(t, "%d Tage", translation)
}

func TestGetCLDRPluralFunction(t *testing.T) {
	pf := getCLDRPluralFunction(language.MustParse("kw"))

	tests := []struct {
		input int
		ouput int
	}{
		{0, 0},
		{1, 1},
		{22, 2},
		{143, 3},
		{161, 4},
		{5, 5},
		{1004, 5},
	}

	for _, test := range tests {
		form, err := pf(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.ouput, form)
	}
}
