package spreak

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/localize"
)

func getLocale(t *testing.T) *Locale {
	bundle, errB := NewBundle(
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	locale, errL := NewLocale(bundle, language.German)
	require.NoError(t, errL)
	require.NotNil(t, locale)

	require.Equal(t, "a", locale.defaultDomain)

	return locale
}

func TestLocale_Getter(t *testing.T) {
	locale := getLocale(t)

	assert.Equal(t, language.German, locale.Language())
	assert.Equal(t, "a", locale.DefaultDomain())
	assert.Equal(t, "z", locale.WithDomain("z").DefaultDomain())
	assert.ElementsMatch(t, []string{"a", "z"}, locale.Domains())
	assert.True(t, locale.HasDomain("a"))
	assert.True(t, locale.HasDomain("z"))
	assert.False(t, locale.HasDomain("b"))
}

func TestLocale_SimplePublicFunctions(t *testing.T) {
	locale := getLocale(t)
	assert.Equal(t, "ID", locale.Get("id"))
	assert.Equal(t, "en_id", locale.Get("en_id"))
	assert.Equal(t, "en_id", locale.Getf("en_id"))

	assert.Equal(t, "ID", locale.DGet("a", "id"))
	assert.Equal(t, "en_id", locale.DGet("a", "en_id"))
	assert.Equal(t, "id", locale.DGet("unknown", "id"))
	assert.Equal(t, "id", locale.DGetf("unknown", "id"))

	assert.Equal(t, "%d Tag", locale.NGet("%d day", "%d days", 1))
	assert.Equal(t, "1 Tag", locale.NGetf("%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.NGetf("%d car", "%d cars", 10, 10))

	assert.Equal(t, "%d Tag", locale.DNGet("a", "%d day", "%d days", 1))
	assert.Equal(t, "1 Tag", locale.DNGetf("a", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.DNGetf("a", "%d car", "%d cars", 10, 10))
	assert.Equal(t, "1 day", locale.DNGetf("unknown", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.DNGetf("unknown", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "Test mit Context", locale.PGet("context", "Test with context"))
	assert.Equal(t, "Test with context", locale.PGetf("other", "Test with context", 5))

	assert.Equal(t, "Test mit Context", locale.DPGet("a", "context", "Test with context"))
	assert.Equal(t, "Test mit Context", locale.DPGetf("a", "context", "Test with context"))
	assert.Equal(t, "Test with context", locale.DPGetf("a", "other", "Test with context", 5))
	assert.Equal(t, "Test with context", locale.DPGetf("unknown", "context", "Test with context"))
	assert.Equal(t, "Test with context", locale.DPGetf("unknown", "other", "Test with context", 5))

	assert.Equal(t, "1 Ergebniss", locale.NPGetf("context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 Ergebnisse", locale.NPGetf("context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", locale.NPGetf("other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", locale.NPGetf("other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", locale.NPGetf("context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", locale.NPGetf("context", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 Ergebniss", locale.DNPGetf("a", "context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 Ergebnisse", locale.DNPGetf("a", "context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", locale.DNPGetf("a", "other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", locale.DNPGetf("a", "other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", locale.DNPGetf("a", "context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", locale.DNPGetf("a", "context", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 result", locale.DNPGetf("unknown", "context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", locale.DNPGetf("unknown", "context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", locale.DNPGetf("unknown", "other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", locale.DNPGetf("unknown", "other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", locale.DNPGetf("unknown", "context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", locale.DNPGetf("unknown", "context", "%d car", "%d cars", 10, 10))
}

func TestLocale_TranslateWithError(t *testing.T) {
	locale := getLocale(t)

	localizeMsg := &localize.Message{
		Singular: "%d day",
		Plural:   "%d days",
		Context:  "",
		Vars:     []interface{}{10},
		Count:    10,
	}

	tr, err := locale.LocalizeWithError(localizeMsg)
	assert.NoError(t, err)
	assert.Equal(t, "10 Tage", tr)

	localizeMsg.Count = 1
	localizeMsg.Vars[0] = 1
	assert.Equal(t, "1 Tag", locale.Localize(localizeMsg))

	localizeMsg.Singular = "Test with context"
	localizeMsg.Plural = ""
	localizeMsg.Count = 0
	localizeMsg.Context = "context"
	tr, err = locale.LocalizeWithError(localizeMsg)
	assert.NoError(t, err)
	assert.Equal(t, "Test mit Context", tr)
}

func TestLocale_Localize(t *testing.T) {
	locale := getLocale(t)
	msg := &localize.Message{
		Singular: "%d day",
		Plural:   "%d days",
		Context:  "",
		Vars:     nil,
		Count:    10,
	}

	tr, err := locale.LocalizeWithError(msg)
	assert.NoError(t, err)
	assert.Equal(t, "%d Tage", tr)

	msg.Count = 0
	tr, err = locale.LocalizeWithError(msg)
	assert.NoError(t, err)
	assert.Equal(t, "%d Tage", tr)
}

func TestLocale_MainFunctions(t *testing.T) {
	locale := getLocale(t)

	t.Run("translate singular", func(t *testing.T) {
		for _, tt := range singularTestData {
			got, err := locale.dpGettextErr("a", tt.ctx, tt.msgID, tt.params...)
			tt.wantErr(t, err, fmt.Sprintf("pGettextErr(%q, %v)", tt.ctx, tt.msgID))
			assert.Equalf(t, tt.translated, got, "pGettextErr(%q, %v)", tt.ctx, tt.msgID)
		}
	})

	t.Run("translate plural", func(t *testing.T) {
		for idx, tt := range pluralTestData {
			got, err := locale.dnpGettextErr("a", tt.ctx, tt.msgID, tt.plural, tt.n, tt.params...)
			if tt.wantErr(t, err, fmt.Sprintf("locale idx=%d dnpGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)) {
				assert.Equalf(t, tt.translated, got, "locale idx=%d dnpGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)
			}
		}
	})

	t.Run("default text is returned for a non-existent domain", func(t *testing.T) {
		tr, err := locale.dnpGettextErr("unknown", NoCtx, "%d day", "%d days", 1, 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)

		tr, err = locale.dnpGettextErr("unknown", NoCtx, "%d day", "%d days", 10, 10)
		assert.Error(t, err)
		assert.Equal(t, "10 days", tr)

		tr, err = locale.dpGettextErr("unknown", NoCtx, "%d day", 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)
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

	locale, errL := NewLocale(bundle, language.Italian)
	require.NoError(t, errL)
	require.NotNil(t, locale)
	assert.Equal(t, language.Italian, locale.language)
	assert.True(t, locale.isSourceLanguage)
}

func TestLocale_LocalizeError(t *testing.T) {
	locale := getLocale(t)

	t.Run("error translation", func(t *testing.T) {
		want := "fehler"
		get := locale.LocalizeError(errors.New("failure"))
		assert.Error(t, get)
		assert.IsType(t, &localize.Error{}, get)
		assert.Equal(t, want, get.Error())
	})

	t.Run("localizable error", func(t *testing.T) {
		raw := &testLocalizeErr{
			singular:  "failure",
			errorText: "other",
		}

		get := locale.LocalizeError(raw)
		assert.Error(t, get)
		assert.IsType(t, raw, get)
		assert.Equal(t, raw.errorText, get.Error())

		raw.context = "errors"
		get = locale.LocalizeError(raw)
		assert.Error(t, get)
		assert.IsType(t, &localize.Error{}, get)
		want := "fehler"
		assert.Equal(t, want, get.Error())

		raw.singular = "The kindergarten"
		raw.context = ""
		raw.domain = "z"
		raw.hasDomain = true
		get = locale.LocalizeError(raw)
		assert.Error(t, get)
		if assert.IsType(t, &localize.Error{}, get) {
			loErr := get.(*localize.Error)
			want = "Der Kindergarten"
			assert.Equal(t, want, loErr.Translation)
		}
	})

	t.Run("no translation", func(t *testing.T) {
		err := errors.New("unknown error")
		get := locale.LocalizeError(err)
		assert.Exactly(t, err, get)
	})
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

	locale, errL := NewLocale(bundle, language.German)
	require.NoError(t, errL)
	require.NotNil(t, locale)

	want := "empty translation"
	get := locale.Get(want)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingTranslation{}, lastErr)

	get = locale.NGet(want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingTranslation{}, lastErr)

	get = locale.PGet("unknown", want)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingContext{}, lastErr)

	get = locale.NPGet("unknown", want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingContext{}, lastErr)

	get = locale.DGet("unknown-domain", want)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingDomain{}, lastErr)

	get = locale.DNGet("unknown-domain", want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingDomain{}, lastErr)

	want = "#+/?%=§$=§$"
	get = locale.Get(want)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingMessageID{}, lastErr)

	get = locale.NGet(want, "plural", 1)
	assert.Equal(t, want, get)
	assert.IsType(t, &ErrMissingMessageID{}, lastErr)
}
