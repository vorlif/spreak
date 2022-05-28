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
