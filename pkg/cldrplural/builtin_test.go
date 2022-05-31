package cldrplural

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestForLanguage(t *testing.T) {
	t.Run("if the language is not known a fallback is used", func(t *testing.T) {
		unknownLang := language.MustParse("x-art")
		forms, found := ForLanguage(unknownLang)
		require.False(t, found)
		require.NotNil(t, forms)

		fallbackForm := builtInRuleSets[language.English.String()]
		require.NotNil(t, fallbackForm)

		sf1 := reflect.ValueOf(fallbackForm).Pointer()
		sf2 := reflect.ValueOf(forms).Pointer()
		assert.Equal(t, sf1, sf2)
	})

	t.Run("if the language has a region the base is used as fallback", func(t *testing.T) {
		ar := language.MustParse("ar")
		arForms, arFound := ForLanguage(ar)
		require.True(t, arFound)
		require.NotNil(t, arForms)

		arDe := language.MustParse("ar-DE")
		deForms, deFound := ForLanguage(arDe)
		require.True(t, deFound)
		require.NotNil(t, deForms)

		sf1 := reflect.ValueOf(arForms).Pointer()
		sf2 := reflect.ValueOf(deForms).Pointer()
		assert.Equal(t, sf1, sf2)
	})

	t.Run("returns the correct value", func(t *testing.T) {
		tests := []string{"de", "en", "dz", "ar", "uk", "zh", "lag"}
		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				lang := language.MustParse(tt)
				got, gotFound := ForLanguage(lang)
				forms := builtInRuleSets[lang.String()]
				assert.Equalf(t, forms, got, "ForLanguage(%v)", lang)
				assert.True(t, gotFound, "ForLanguage(%v)", lang)
			})
		}
	})

	t.Run("uses backward compatibility", func(t *testing.T) {
		edgeTests := []struct {
			langName string
			wantLang string
		}{
			{"nl-BE", "nl"},
			{"zh-Hant", "zh"},
			{"de-AT", "de"},
		}
		for _, tt := range edgeTests {
			t.Run(tt.langName, func(t *testing.T) {
				lang := language.MustParse(tt.langName)
				got, gotFound := ForLanguage(lang)
				forms := builtInRuleSets[tt.wantLang]
				assert.Equalf(t, forms, got, "pluralRuleForLanguage(%v)", lang)
				assert.True(t, gotFound, "pluralRuleForLanguage(%v)", lang)
			})
		}
	})
}
