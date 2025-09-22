package poplural

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
	// github.com/vorlif/spreak/catalog/poplural.forLanguage.newFormCsSk.func28
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		// github.com/vorlif/spreak/catalog/poplural.forLanguage.newFormCsSk
		return funcName[:idx]
	}
	return funcName
}

func TestForLanguage(t *testing.T) {
	t.Run("if the language is not known a fallback is used", func(t *testing.T) {
		unknownLang := language.MustParse("art")
		rule, found := ForLanguage(unknownLang)
		require.False(t, found)
		require.NotNil(t, rule)

		fallbackRule := getBuiltInForLanguage(fallbackLanguage)
		require.NotNil(t, fallbackRule)
		assert.Equal(t, getFuncName(fallbackRule.Evaluate), getFuncName(rule))
	})

	t.Run("if the language has a region the base is used as fallback", func(t *testing.T) {
		ar := language.MustParse("ar")
		arRule, arFound := ForLanguage(ar)
		require.True(t, arFound)
		require.NotNil(t, arRule)

		arDe := language.MustParse("ar-DE")
		deRule, deFound := ForLanguage(arDe)
		require.True(t, deFound)
		require.NotNil(t, deRule)

		assert.Equal(t, getFuncName(arRule), getFuncName(deRule))
	})
}

func Test_pluralRuleForLanguage(t *testing.T) {
	tests := []string{"de", "en", "dz", "ar", "uk", "zh", "lag"}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			lang := language.MustParse(tt)
			got, gotFound := pluralRuleForLanguage(lang)
			forms := getBuiltInForLanguage(lang.String())
			assert.Equal(t, forms.NPlurals, got.NPlurals, "pluralRuleForLanguage(%v).NPlurals", lang)
			assert.Equal(t, getFuncName(forms.FormFunc), getFuncName(got.FormFunc), "pluralRuleForLanguage(%v).FormFunc", lang)
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
			forms := getBuiltInForLanguage(tt.wantLang)
			assert.Equal(t, forms.NPlurals, got.NPlurals, "pluralRuleForLanguage(%v).NPlurals", lang)
			assert.Equal(t, getFuncName(forms.FormFunc), getFuncName(got.FormFunc), "pluralRuleForLanguage(%v).FormFunc", lang)
			assert.True(t, gotFound, "pluralRuleForLanguage(%v)", lang)
		})
	}
}
