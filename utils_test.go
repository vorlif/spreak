package spreak

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog"
	"github.com/vorlif/spreak/localize"
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
	want := []string{"de-Latn", "de_Latn", "de_AT", "de-AT", "deu", "de"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.Chinese)
	want = []string{"zh-Hans", "zh_Hans", "zh-CN", "zh_CN", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.TraditionalChinese)
	want = []string{"zh_Hant", "zh-Hant", "zh-TW", "zh_TW", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.SimplifiedChinese)
	want = []string{"zh_Hans", "zh-Hans", "zh-CN", "zh_CN", "zho", "zh"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("sr_LATN"))
	want = []string{"sr_Latn", "sr-Latn", "sr-RS", "sr_RS", "srp", "sr"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("art-x-a2"))
	want = []string{"art-x-a2", "art"}
	assert.Equal(t, want, got)

	got = ExpandLanguage(language.MustParse("en-US-x-twain"))
	want = []string{"en-US-x-twain", "en-Latn", "en_Latn", "en-US", "en_US", "eng", "en"}
	assert.Equal(t, want, got)
}

func TestErrMissingDomain(t *testing.T) {
	err := &ErrMissingDomain{}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spreak")

	err.Language = language.Afrikaans
	assert.Contains(t, err.Error(), err.Language.String())

	err.Domain = "mydomain"
	assert.Contains(t, err.Error(), err.Domain)
}

var (
	testdataStructureDir = filepath.FromSlash("./testdata/structure")
	testTranslationDir   = filepath.FromSlash("testdata/translation-test/")
)

type testDecoder struct {
	f func(lang language.Tag, domain string, data []byte) (catalog.Catalog, error)
}

func (t *testDecoder) Decode(lang language.Tag, domain string, data []byte) (catalog.Catalog, error) {
	return t.f(lang, domain, data)
}

var _ catalog.Decoder = (*testDecoder)(nil)

type testLoader struct {
	f func(lang language.Tag, domain string) (catalog.Catalog, error)
}

var _ Loader = (*testLoader)(nil)

func (t *testLoader) Load(lang language.Tag, domain string) (catalog.Catalog, error) {
	return t.f(lang, domain)
}

type testResolver struct {
	f func(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error)
}

func (t *testResolver) Resolve(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error) {
	return t.f(fsys, extension, lang, domain)
}

var _ Resolver = (*testResolver)(nil)

type testFs struct {
	f func(name string) (fs.File, error)
}

var _ fs.FS = (*testFs)(nil)

func (t *testFs) Open(name string) (fs.File, error) {
	return t.f(name)
}

type testPrinter struct {
	f func(tag language.Tag) PrintFunc
}

func (t *testPrinter) Init(_ []language.Tag) {}

var _ Printer = (*testPrinter)(nil)

func (t *testPrinter) GetPrintFunc(tag language.Tag) PrintFunc { return t.f(tag) }

type testLocalizeErr struct {
	singular  string
	plural    string
	context   string
	domain    string
	hasDomain bool
	errorText string
}

var _ localize.Localizable = (*testLocalizeErr)(nil)

func (t *testLocalizeErr) GetMsgID() string { return t.singular }

func (t *testLocalizeErr) GetPluralID() string { return t.plural }

func (t *testLocalizeErr) GetContext() string { return t.context }

func (testLocalizeErr) GetVars() []any { return nil }

func (testLocalizeErr) GetCount() int { return 0 }

func (t *testLocalizeErr) HasDomain() bool { return t.hasDomain }

func (t *testLocalizeErr) GetDomain() string { return t.domain }

func (t *testLocalizeErr) Error() string { return t.errorText }

func (t *testLocalizeErr) String() string { return t.GetMsgID() }

var singularTestData = []struct {
	msgID      string
	ctx        string
	translated string
	params     []any
	wantErr    assert.ErrorAssertionFunc
}{
	{
		"id",
		NoCtx,
		"ID",
		[]interface{}{},
		assert.NoError,
	},
	{
		"Test with special characters “%s”",
		NoCtx,
		"Test mit Sonderzeichen “abc”",
		[]interface{}{"abc"},
		assert.NoError,
	},
	{
		"Test with context",
		"context",
		"Test mit Context",
		[]interface{}{},
		assert.NoError,
	},
	{
		"unknown text",
		"",
		"unknown text",
		[]interface{}{},
		assert.Error,
	},
	{
		"id",
		"context",
		"id",
		[]interface{}{},
		assert.Error,
	},
	{
		"id",
		"unknown_context",
		"id",
		[]interface{}{},
		assert.Error,
	},
	{
		"empty translation",
		NoCtx,
		"empty translation",
		[]interface{}{},
		assert.Error,
	},
}

var pluralTestData = []struct {
	msgID      string
	plural     string
	n          int
	ctx        string
	translated string
	params     []interface{}
	wantErr    assert.ErrorAssertionFunc
}{
	{
		"%d day",
		"%d days",
		1,
		NoCtx,
		"1 Tag",
		[]interface{}{1},
		assert.NoError,
	},
	{
		"%d day",
		"%d days",
		10,
		NoCtx,
		"10 Tage",
		[]interface{}{10},
		assert.NoError,
	},
	{
		"%d result",
		"%d results",
		1,
		"context",
		"1 Ergebniss",
		[]interface{}{1},
		assert.NoError,
	},
	{
		"%d result",
		"%d results",
		10,
		"context",
		"10 Ergebnisse",
		[]interface{}{10},
		assert.NoError,
	},
	{
		"unknown text",
		"unknown texts",
		1,
		NoCtx,
		"unknown text",
		[]interface{}{},
		assert.Error,
	},
	{
		"unknown text",
		"unknown texts",
		10,
		NoCtx,
		"unknown texts",
		[]interface{}{},
		assert.Error,
	},
	{
		"%d day",
		"%d days",
		10,
		"context",
		"10 days",
		[]interface{}{10},
		assert.Error,
	},
}
