package spreak

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func Test_parseLanguageName(t *testing.T) {

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

func Test_defaultPrintFuncGenerator(t *testing.T) {
	engP := NewDefaultPrinter().GetPrintFunc(language.English)
	require.NotNil(t, engP)

	deP := NewDefaultPrinter().GetPrintFunc(language.German)
	require.NotNil(t, deP)

	test := "test for number %d"
	assert.Equal(t, "test for number 12,345,678", engP(test, 12345678))
	assert.Equal(t, "test for number 12.345.678", deP(test, 12345678))
	assert.Equal(t, test, deP(test))
}
