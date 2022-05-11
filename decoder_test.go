package spreak

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func TestMoDecoder(t *testing.T) {
	dec := NewMoDecoder()
	decode, err := dec.Decode(language.German, "", []byte{})
	assert.Error(t, err)
	assert.Nil(t, decode)
}
