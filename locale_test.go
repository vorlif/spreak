package spreak

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog"
)

func getLocale(t *testing.T) *locale {
	bundle, errB := NewBundle(
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, language.German)
	require.NotNil(t, localizer)

	return localizer.locale
}

func TestLocale_MainFunctions(t *testing.T) {
	tLocale := getLocale(t)

	t.Run("translate singular", func(t *testing.T) {
		for _, tt := range singularTestData {
			got, err := tLocale.lookupSingularTranslation("a", tt.ctx, tt.msgID, tt.params...)
			tt.wantErr(t, err, fmt.Sprintf("lookupSingularTranslation(%q, %v)", tt.ctx, tt.msgID))
			if err != nil {
				assert.Equalf(t, "", got, "lookupSingularTranslation(%q, %v)", tt.ctx, tt.msgID)
			} else {
				assert.Equalf(t, tt.translated, got, "lookupSingularTranslation(%q, %v)", tt.ctx, tt.msgID)
			}
		}
	})

	t.Run("translate plural", func(t *testing.T) {
		for idx, tt := range pluralTestData {
			got, err := tLocale.lookupPluralTranslation("a", tt.ctx, tt.msgID, tt.plural, tt.n, tt.params...)
			tt.wantErr(t, err, fmt.Sprintf("tLocale idx=%d lookupPluralTranslation(%q, %q)", idx, tt.ctx, tt.msgID))
			if err != nil {
				assert.Equalf(t, "", got, "tLocale idx=%d lookupPluralTranslation(%q, %q)", idx, tt.ctx, tt.msgID)
			} else {
				assert.Equalf(t, tt.translated, got, "tLocale idx=%d lookupPluralTranslation(%q, %q)", idx, tt.ctx, tt.msgID)
			}
		}
	})

	t.Run("empty text is returned for a non-existent domain", func(t *testing.T) {
		tr, err := tLocale.lookupPluralTranslation("unknown", NoCtx, "%d day", "%d days", 1, 1)
		assert.Error(t, err)
		assert.Equal(t, "", tr)

		tr, err = tLocale.lookupPluralTranslation("unknown", NoCtx, "%d day", "%d days", 10, 10)
		assert.Error(t, err)
		assert.Equal(t, "", tr)

		tr, err = tLocale.lookupSingularTranslation("unknown", NoCtx, "%d day", 1)
		assert.Error(t, err)
		assert.Equal(t, "", tr)
	})
}

func TestNewLocale_UseSourceLanguage(t *testing.T) {
	bundle, errB := NewBundle(
		WithSourceLanguage(language.Italian),
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	locale := buildSourceLocale(bundle, language.Italian)
	require.NotNil(t, locale)
	assert.Equal(t, language.Italian, locale.language)
	assert.True(t, locale.isSourceLanguage)
}

func TestLocaleMissingCallback(t *testing.T) {
	var lastErr error
	bundle, errB := NewBundle(
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithLanguage(language.German),
		WithMissingTranslationCallback(func(err error) {
			lastErr = err
		}),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	locale := NewLocalizer(bundle, language.German)
	require.NotNil(t, locale)

	want := "empty translation"
	get := locale.Get(want)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingTranslation{}, lastErr)

	get = locale.NGet(want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingTranslation{}, lastErr)

	get = locale.PGet("unknown", want)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingContext{}, lastErr)

	get = locale.NPGet("unknown", want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingContext{}, lastErr)

	get = locale.DGet("unknown-domain", want)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingDomain{}, lastErr)

	get = locale.DNGet("unknown-domain", want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingDomain{}, lastErr)

	want = "#+/?%=§$=§$"
	get = locale.Get(want)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingMessageID{}, lastErr)

	get = locale.NGet(want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &catalog.ErrMissingMessageID{}, lastErr)
}
