package spreak

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/localize"
)

func getLocaleForDomain(t *testing.T) *Locale {
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
	locale.Domains()

	return locale
}

func TestLocale_Getter(t *testing.T) {
	locale := getLocaleForDomain(t)

	assert.Equal(t, language.German, locale.Language())
	assert.Equal(t, "a", locale.DefaultDomain())
	assert.Equal(t, "z", locale.WithDomain("z").DefaultDomain())
	assert.ElementsMatch(t, []string{"a", "z"}, locale.Domains())
	assert.True(t, locale.HasDomain("a"))
	assert.True(t, locale.HasDomain("z"))
	assert.False(t, locale.HasDomain("b"))
}

func TestLocale_SimplePublicFunctions(t *testing.T) {
	locale := getLocaleForDomain(t)
	assert.Equal(t, "ID", locale.Get("id"))
	assert.Equal(t, "en_id", locale.Get("en_id"))

	assert.Equal(t, "ID", locale.DGet("a", "id"))
	assert.Equal(t, "en_id", locale.DGet("a", "en_id"))
	assert.Equal(t, "id", locale.DGet("unknown", "id"))

	assert.Equal(t, "1 Tag", locale.NGetf("%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.NGetf("%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 Tag", locale.DNGetf("a", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.DNGetf("a", "%d car", "%d cars", 10, 10))
	assert.Equal(t, "1 day", locale.DNGetf("unknown", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", locale.DNGetf("unknown", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "Test mit Context", locale.PGet("context", "Test with context"))
	assert.Equal(t, "Test with context", locale.PGetf("other", "Test with context", 5))

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
	locale := getLocaleForDomain(t)

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

func TestLocale_MainFunctions(t *testing.T) {
	locale := getLocaleForDomain(t)

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
