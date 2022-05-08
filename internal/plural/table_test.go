package plural

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestPluralFormsTable(t *testing.T) {
	for lang, rule := range pluralFormsTable {
		_, err := Parse(rule)
		if err != nil {
			t.Errorf("Parse(%s) for %s generates error", rule, lang)
		}
	}
}

func TestForLanguage(t *testing.T) {
	t.Run("if the language is not known a fallback is used", func(t *testing.T) {
		unknownLang := language.MustParse("art")
		forms, found := ForLanguage(unknownLang)
		require.False(t, found)
		require.NotNil(t, forms)

		fallbackRule := pluralFormsTable[language.Und.String()]
		fallbackForms, err := Parse(fallbackRule)
		require.NoError(t, err)
		require.NotNil(t, fallbackForms)

		sf1 := reflect.ValueOf(fallbackForms.IndexForN).Pointer()
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
}

func Test_pluralRuleForLanguage(t *testing.T) {
	tests := []string{"de", "en", "dz", "ar", "uk", "zh", "lag"}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			lang := language.MustParse(tt)
			got, gotFound := pluralRuleForLanguage(lang)
			forms := pluralFormsTable[tt]
			assert.Equalf(t, forms, got, "pluralRuleForLanguage(%v)", lang)
			assert.True(t, gotFound, "pluralRuleForLanguage(%v)", lang)
		})
	}

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
			got, gotFound := pluralRuleForLanguage(lang)
			forms := pluralFormsTable[tt.wantLang]
			assert.Equalf(t, forms, got, "pluralRuleForLanguage(%v)", lang)
			assert.True(t, gotFound, "pluralRuleForLanguage(%v)", lang)
		})
	}
}
