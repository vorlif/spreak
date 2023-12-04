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

	translation, err := cat.Lookup("", "id")
	if assert.NoError(t, err) {
		assert.Equal(t, "ID", translation)
	}

	translation, err = cat.LookupPlural("", "%d day", 1)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Tag", translation)
	}

	translation, err = cat.LookupPlural("", "%d car", 10)
	if assert.Error(t, err) {
		assert.Equal(t, "%d car", translation)
	}

	translation, err = cat.Lookup("context", "Test with context")
	if assert.NoError(t, err) {
		assert.Equal(t, "Test mit Context", translation)
	}

	translation, err = cat.LookupPlural("other", "Test with context", 5)
	if assert.Error(t, err) {
		assert.Equal(t, "Test with context", translation)
	}

	translation, err = cat.LookupPlural("context", "%d result", 1)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Ergebniss", translation)
	}

	translation, err = cat.LookupPlural("context", "%d result", 10)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Ergebnisse", translation)
	}

	translation, err = cat.LookupPlural("other", "%d result", 1)
	if assert.Error(t, err) {
		assert.Equal(t, "%d result", translation)
	}

	translation, err = cat.LookupPlural("other", "%d result", 10)
	if assert.Error(t, err) {
		assert.Equal(t, "%d result", translation)
	}
}

func TestGetCLDRPluralFunction(t *testing.T) {
	pf := getCLDRPluralFunction(language.MustParse("kw"))

	tests := []struct {
		input  int
		output int
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
		assert.Equal(t, test.output, form)
	}
}

func TestGettextMessage_Translations(t *testing.T) {
	msg := GettextMessage{
		translations: map[int]string{0: "Car", 1: "Cars"},
	}

	cpy := msg.Translations()
	assert.Equal(t, msg.translations, cpy)
	cpy[0] = "Auto"
	cpy[1] = "Autos"
	assert.NotEqual(t, msg.translations, cpy)
}
