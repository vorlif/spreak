package cldrplural

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// Returns the name of a function for a function type.
func getFuncName(i any) string {
	// github.com/vorlif/spreak/v2/catalog/poplural.forLanguage.newFormCsSk.func28
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		// github.com/vorlif/spreak/v2/catalog/poplural.forLanguage.newFormCsSk
		return funcName[:idx]
	}
	return funcName
}

func TestForLanguage(t *testing.T) {
	t.Run("if the language is not known a fallback is used", func(t *testing.T) {
		unknownLang := language.MustParse("x-art")
		forms, found := ForLanguage(unknownLang)
		require.False(t, found)
		require.NotNil(t, forms)

		fallbackForm := getBuiltInForLanguage(fallbackLanguage)
		require.NotNil(t, fallbackForm)
		assert.Equal(t, getFuncName(fallbackForm), getFuncName(forms))
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

		assert.Equal(t, getFuncName(arForms), getFuncName(deForms))
	})

	t.Run("returns the correct value", func(t *testing.T) {
		tests := []string{"de", "en", "dz", "ar", "uk", "zh", "lag"}
		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				lang := language.MustParse(tt)
				got, gotFound := ForLanguage(lang)
				forms := getBuiltInForLanguage(lang.String())
				assert.Equalf(t, getFuncName(forms), getFuncName(got), "ForLanguage(%v)", lang)
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
				forms := getBuiltInForLanguage(lang.String())
				assert.Equalf(t, getFuncName(forms), getFuncName(got), "ForLanguage(%v)", lang)
				assert.True(t, gotFound, "ForLanguage(%v)", lang)
			})
		}
	})
}
