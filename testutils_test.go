package spreak

import (
	"io/fs"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/localize"
)

var (
	testdataStructureDir = filepath.FromSlash("./testdata/structure")
	testTranslationDir   = filepath.FromSlash("testdata/translation-test/")
)

type testDecoder struct {
	f func(lang language.Tag, domain string, data []byte) (Catalog, error)
}

func (t *testDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	return t.f(lang, domain, data)
}

var _ Decoder = (*testDecoder)(nil)

type testLoader struct {
	f func(lang language.Tag, domain string) (Catalog, error)
}

var _ Loader = (*testLoader)(nil)

func (t *testLoader) Load(lang language.Tag, domain string) (Catalog, error) {
	return t.f(lang, domain)
}

type testReducer struct {
	f func(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error)
}

func (t *testReducer) Reduce(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error) {
	return t.f(fsys, extension, lang, domain)
}

var _ Reducer = (*testReducer)(nil)

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

func (testLocalizeErr) GetVars() []interface{} { return nil }

func (testLocalizeErr) GetCount() int { return 0 }

func (t *testLocalizeErr) HasDomain() bool { return t.hasDomain }

func (t *testLocalizeErr) GetDomain() string { return t.domain }

func (t *testLocalizeErr) Error() string { return t.errorText }

func (t *testLocalizeErr) String() string { return t.GetMsgID() }

var singularTestData = []struct {
	msgID      string
	ctx        string
	translated string
	params     []interface{}
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
