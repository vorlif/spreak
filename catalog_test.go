package spreak

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func getCatalogForDomainA(t *testing.T, domain string) Catalog {
	bundle, errB := NewBundle(
		WithDefaultDomain(domain),
		WithDomainPath(domain, testTranslationDir),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	locale, errL := NewLocale(bundle, language.German)
	require.NoError(t, errL)
	require.NotNil(t, locale)

	cat, errC := locale.getCatalog(locale.defaultDomain)
	require.NoError(t, errC)
	require.NotNil(t, cat)
	return cat
}

func TestCatalog_SimplePublicFunctions(t *testing.T) {
	cat := getCatalogForDomainA(t, "a")
	assert.Equal(t, language.German, cat.Language())

	translation, err := cat.GetTranslation(NoCtx, "id")
	if assert.NoError(t, err) {
		assert.Equal(t, "ID", translation)
	}

	translation, err = cat.GetPluralTranslation(NoCtx, "%d day", 1)
	if assert.NoError(t, err) {
		assert.Equal(t, "%d Tag", translation)
	}

	translation, err = cat.GetPluralTranslation(NoCtx, "%d car", 10)
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
