package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

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

#, fuzzy
msgid "Unknown"
msgstr "Fuzzy translation"
`

func TestPoDecoder(t *testing.T) {
	domain := "my-domain"
	lang := language.German

	dec := NewPoDecoder()
	cat, err := dec.Decode(lang, domain, []byte{})
	assert.Error(t, err)
	assert.Nil(t, cat)

	cat, err = dec.Decode(lang, domain, []byte(decodeTestData))
	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, lang, cat.Language())

	translation, err := cat.Lookup("", "id")
	assert.NoError(t, err)
	assert.Equal(t, "ID", translation)

	translation, err = cat.Lookup("", "xyz")
	assert.Error(t, err)
	assert.Equal(t, "xyz", translation)

	translation, err = cat.Lookup("", "Unknown")
	assert.Error(t, err)
	assert.Equal(t, "Unknown", translation)

	cat, err = dec.Decode(lang, domain, []byte(decodePoInvalidPluralForm+decodeTestData))
	assert.Error(t, err)
	assert.Nil(t, cat)

	t.Run("test special chars", func(t *testing.T) {
		cat, err = dec.Decode(lang, domain, []byte(decodeTestData))
		assert.NoError(t, err)
		require.NotNil(t, cat)

		tr, err := cat.Lookup("", "Special		\"chars\"")
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

	cat, err := dec.Decode(language.German, "", []byte(decodeTestData))
	assert.NoError(t, err)
	require.NotNil(t, cat)

	translation, errT := cat.LookupPlural("", "%d day", "1.2")
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
	cat, err := dec.Decode(language.German, "", []byte(header+decodeTestData))
	assert.NoError(t, err)
	require.NotNil(t, cat)
	translation, errT := cat.LookupPlural("", "%d day", "1.2")
	assert.NoError(t, errT)
	assert.Equal(t, "%d Tage", translation)
}
