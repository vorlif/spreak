package spreak

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewBundle(t *testing.T) {
	t.Run("error is returned when a nil option is passed", func(t *testing.T) {
		bundle, err := NewBundle(WithDomainPath(NoDomain, testdataStructureDir), nil)
		require.Error(t, err)
		require.Nil(t, bundle)
	})
}

func TestBundle_Domains(t *testing.T) {
	bundle, err := NewBundle(
		WithDomainPath(NoDomain, testdataStructureDir),
		WithDomainPath("a", testdataStructureDir),
		WithDomainPath("b", testdataStructureDir),
	)
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, 0, len(bundle.Domains()))
	assert.Equal(t, 0, len(bundle.SupportedLanguages()))
}

func TestBundle_SupportedLanguages(t *testing.T) {
	bundle, err := NewBundle(
		WithDomainPath(NoDomain, testdataStructureDir),
		WithRequiredLanguage(language.English, language.MustParse("de-at")),
		WithLanguage(language.Afrikaans),
	)
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, 1, len(bundle.Domains()))
	assert.Equal(t, 2, len(bundle.SupportedLanguages()))
}

func TestBundle_GetLocaleWithDomain(t *testing.T) {
	reducer, errRed := NewDefaultResolver(WithCategory("my_category"))
	require.NoError(t, errRed)
	loader, errL := NewFilesystemLoader(WithResolver(reducer), WithPath(testdataStructureDir))
	require.NoError(t, errL)
	require.NotNil(t, loader)

	bundle, errB := NewBundle(
		WithDomainPath(NoDomain, testdataStructureDir),
		WithDomainPath("a", testdataStructureDir),
		WithDomainLoader("b", loader),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	t.Run("if the language has not been loaded, no locale can be created", func(t *testing.T) {
		locale, err := NewLocaleWithDomain(bundle, language.Danish, NoDomain)
		assert.Error(t, err)
		assert.Nil(t, locale)
	})

	t.Run("when the language is loaded, the locale can be created", func(t *testing.T) {
		locale, err := NewLocaleWithDomain(bundle, language.German, NoDomain)
		assert.NoError(t, err)
		assert.NotNil(t, locale)
		assert.Equal(t, language.German, locale.Language())
		assert.Equal(t, NoDomain, locale.DefaultDomain())
		assert.ElementsMatch(t, []string{"a", "b"}, locale.Domains())
	})

	t.Run("another domain returns a customized locale", func(t *testing.T) {
		locale, err := NewLocaleWithDomain(bundle, language.German, "a")
		require.NoError(t, err)
		require.NotNil(t, loader)
		assert.Equal(t, "a", locale.defaultDomain)
	})

}
