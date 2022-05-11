package spreak

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestParseLanguageName(t *testing.T) {

	tests := []struct {
		langStr string
		want    language.Tag
		wantErr assert.ErrorAssertionFunc
	}{
		{"de", language.German, assert.NoError},
		{"de_DE.ISO-8859-1", language.MustParse("de_DE"), assert.NoError},
		{"en@quot", language.English, assert.NoError},
		{"sr_RS@latin", language.MustParse("sr_RS"), assert.NoError},
	}
	for _, tt := range tests {
		got, err := parseLanguageName(tt.langStr)
		if !tt.wantErr(t, err, fmt.Sprintf("parseLanguageName(%v)", tt.langStr)) {
			return
		}
		assert.Equalf(t, tt.want, got, "parseLanguageName(%v)", tt.langStr)
	}
}

func TestDefaultPrintFuncGenerator(t *testing.T) {
	engP := NewDefaultPrinter().GetPrintFunc(language.English)
	require.NotNil(t, engP)

	deP := NewDefaultPrinter().GetPrintFunc(language.German)
	require.NotNil(t, deP)

	test := "test for number %d"
	assert.Equal(t, "test for number 12,345,678", engP(test, 12345678))
	assert.Equal(t, "test for number 12.345.678", deP(test, 12345678))
	assert.Equal(t, test, deP(test))
}

func TestExpandLanguage(t *testing.T) {
	got := ExpandLanguage(language.MustParse("de-AT"))
	want := []string{"de-Latn", "de_Latn", "de-AT", "de_AT", "deu", "de"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.Chinese)
	want = []string{"zh-Hans", "zh_Hans", "zh-CN", "zh_CN", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.TraditionalChinese)
	want = []string{"zh-Hant", "zh_Hant", "zh-TW", "zh_TW", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.SimplifiedChinese)
	want = []string{"zh-Hans", "zh_Hans", "zh-CN", "zh_CN", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("sr_LATN"))
	want = []string{"sr-Latn", "sr_Latn", "sr-RS", "sr_RS", "srp", "sr"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("art-x-a2"))
	want = []string{"art-x-a2", "art"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("en-US-x-twain"))
	want = []string{"en-US-x-twain", "en-Latn", "en_Latn", "en-US", "en_US", "eng", "en"}
	assert.Equal(t, want, got)
}
