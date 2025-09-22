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
	catl, err := dec.Decode(lang, domain, []byte{})
	assert.Error(t, err)
	assert.Nil(t, catl)

	catl, err = dec.Decode(lang, domain, []byte(decodeTestData))
	assert.NoError(t, err)
	assert.NotNil(t, catl)
	assert.Equal(t, lang, catl.Language())

	translation, err := catl.Lookup("", "id")
	assert.NoError(t, err)
	assert.Equal(t, "ID", translation)

	translation, err = catl.Lookup("", "xyz")
	assert.Error(t, err)
	assert.Equal(t, "xyz", translation)

	translation, err = catl.Lookup("", "Unknown")
	assert.Error(t, err)
	assert.Equal(t, "Unknown", translation)

	catl, err = dec.Decode(lang, domain, []byte(decodePoInvalidPluralForm+decodeTestData))
	assert.Error(t, err)
	assert.Nil(t, catl)

	t.Run("test special chars", func(t *testing.T) {
		catl, err = dec.Decode(lang, domain, []byte(decodeTestData))
		assert.NoError(t, err)
		require.NotNil(t, catl)

		tr, err := catl.Lookup("", "Special		\"chars\"")
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

	translation, errT := catl.LookupPlural("", "%d day", "1.2")
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
	translation, errT := catl.LookupPlural("", "%d day", "1.2")
	assert.NoError(t, errT)
	assert.Equal(t, "%d Tage", translation)
}
