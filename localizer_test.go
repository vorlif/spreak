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

func getLocalizerForTest(t *testing.T) *Localizer {
	bundle, errB := NewBundle(
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithRequiredLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, "de", "en", "en_US")
	require.NotNil(t, localizer)

	return localizer
}

func TestLocalizer_FallbackLanguageIsUsed(t *testing.T) {
	bundle, errB := NewBundle(
		WithSourceLanguage(language.English),
		WithFallbackLanguage(language.German),
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, "zh")
	require.NotNil(t, localizer)

	assert.Equal(t, language.German, localizer.Language())
}

func TestLocalizer_SourceIsUsed(t *testing.T) {
	bundle, errB := NewBundle(
		WithSourceLanguage(language.English),
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, "zh")
	require.NotNil(t, localizer)

	assert.Equal(t, language.English, localizer.Language())
}

func TestNewLocalizer_InvalidInputIsIgnored(t *testing.T) {
	bundle, errB := NewBundle(
		WithSourceLanguage(language.English),
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, 5, 2, nil, "de")
	require.NotNil(t, localizer)
	assert.Equal(t, language.German, localizer.Language())
}

func TestLocalizer_Getter(t *testing.T) {
	localizer := getLocalizerForTest(t)

	assert.Equal(t, language.German, localizer.Language())
	assert.Equal(t, "a", localizer.DefaultDomain())
	assert.True(t, localizer.HasLocale())
	assert.NotNil(t, localizer.GetLocale())
}

func TestLocalizer_NoLocale(t *testing.T) {
	bundle, errB := NewBundle(
		WithDefaultDomain("a"),
		WithDomainPath("a", testTranslationDir),
		WithDomainPath("z", testTranslationDir),
		WithRequiredLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	localizer := NewLocalizer(bundle, "zh-Hans")
	require.NotNil(t, localizer)

	assert.False(t, localizer.HasLocale())
	tr, err := localizer.dnpGettextErr("unknown", NoCtx, "%d day", "%d days", 10, 10)
	assert.Error(t, err)
	assert.Equal(t, "10 days", tr)

	tr, err = localizer.dpGettextErr("unknown", NoCtx, "%d day", 1)
	assert.Error(t, err)
	assert.Equal(t, "1 day", tr)

	assert.Nil(t, localizer.GetLocale())

	localizeMsg := &localize.Message{
		Singular: "%d day",
		Plural:   "%d days",
		Context:  "",
		Vars:     []interface{}{10},
		Count:    10,
	}
	tr, err = localizer.LocalizeWithError(localizeMsg)
	assert.Error(t, err)
	assert.Equal(t, "10 days", tr)
}

func TestLocalizer_SimplePublicFunctions(t *testing.T) {
	localizer := getLocalizerForTest(t)
	assert.Equal(t, "ID", localizer.Getf("id"))
	assert.Equal(t, "en_id", localizer.Getf("en_id"))

	assert.Equal(t, "ID", localizer.DGetf("a", "id"))
	assert.Equal(t, "en_id", localizer.DGetf("a", "en_id"))
	assert.Equal(t, "id", localizer.DGetf("unknown", "id"))

	assert.Equal(t, "1 Tag", localizer.NGetf("%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.NGetf("%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 Tag", localizer.DNGetf("a", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNGetf("a", "%d car", "%d cars", 10, 10))
	assert.Equal(t, "1 day", localizer.DNGetf("unknown", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNGetf("unknown", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "Test mit Context", localizer.PGet("context", "Test with context"))
	assert.Equal(t, "Test mit Context", localizer.PGetf("context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.PGetf("other", "Test with context", 5))

	assert.Equal(t, "Test mit Context", localizer.DPGetf("a", "context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.DPGetf("a", "other", "Test with context", 5))
	assert.Equal(t, "Test with context", localizer.DPGetf("unknown", "context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.DPGetf("unknown", "other", "Test with context", 5))

	assert.Equal(t, "1 Ergebniss", localizer.NPGetf("context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 Ergebnisse", localizer.NPGetf("context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", localizer.NPGetf("other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", localizer.NPGetf("other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", localizer.NPGetf("context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", localizer.NPGetf("context", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 Ergebniss", localizer.DNPGetf("a", "context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 Ergebnisse", localizer.DNPGetf("a", "context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", localizer.DNPGetf("a", "other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", localizer.DNPGetf("a", "other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", localizer.DNPGetf("a", "context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNPGetf("a", "context", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 result", localizer.DNPGetf("unknown", "context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", localizer.DNPGetf("unknown", "context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", localizer.DNPGetf("unknown", "other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", localizer.DNPGetf("unknown", "other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", localizer.DNPGetf("unknown", "context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNPGetf("unknown", "context", "%d car", "%d cars", 10, 10))
}

func TestLocalizer_TranslateWithError(t *testing.T) {
	localizer := getLocalizerForTest(t)

	localizeMsg := &localize.Message{
		Singular: "%d day",
		Plural:   "%d days",
		Context:  "",
		Vars:     []interface{}{10},
		Count:    10,
	}

	tr, err := localizer.LocalizeWithError(localizeMsg)
	assert.NoError(t, err)
	assert.Equal(t, "10 Tage", tr)

	localizeMsg.Count = 1
	localizeMsg.Vars[0] = 1
	assert.Equal(t, "1 Tag", localizer.Localize(localizeMsg))

	localizeMsg.Singular = "Test with context"
	localizeMsg.Plural = ""
	localizeMsg.Count = 0
	localizeMsg.Context = "context"
	tr, err = localizer.LocalizeWithError(localizeMsg)
	assert.NoError(t, err)
	assert.Equal(t, "Test mit Context", tr)
}

func TestLocalizer_MainFunctions(t *testing.T) {
	localizer := getLocalizerForTest(t)

	t.Run("translate singular", func(t *testing.T) {
		for _, tt := range singularTestData {
			got, err := localizer.dpGettextErr("a", tt.ctx, tt.msgID, tt.params...)
			tt.wantErr(t, err, fmt.Sprintf("pGettextErr(%q, %v)", tt.ctx, tt.msgID))
			assert.Equalf(t, tt.translated, got, "pGettextErr(%q, %v)", tt.ctx, tt.msgID)
		}
	})

	t.Run("translate plural", func(t *testing.T) {
		for idx, tt := range pluralTestData {
			got, err := localizer.dnpGettextErr("a", tt.ctx, tt.msgID, tt.plural, tt.n, tt.params...)
			if tt.wantErr(t, err, fmt.Sprintf("localizer idx=%d npGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)) {
				assert.Equalf(t, tt.translated, got, "localizer idx=%d npGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)
			}
		}
	})

	t.Run("default text is returned for a non-existent domain", func(t *testing.T) {
		tr, err := localizer.dnpGettextErr("unknown", NoCtx, "%d day", "%d days", 1, 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)

		tr, err = localizer.dnpGettextErr("unknown", NoCtx, "%d day", "%d days", 10, 10)
		assert.Error(t, err)
		assert.Equal(t, "10 days", tr)

		tr, err = localizer.dpGettextErr("unknown", NoCtx, "%d day", 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)
	})
}

func TestLocalizer_LocalizeError(t *testing.T) {
	localizer := getLocalizerForTest(t)

	err := errors.New("test error")
	err = localizer.LocalizeError(err)
	assert.Error(t, err)
}
