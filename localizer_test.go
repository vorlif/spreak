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
	tr, err := localizer.lookupPluralTranslation("unknown", NoCtx, "%d day", "%d days", 10, 10)
	assert.Error(t, err)
	assert.Equal(t, "10 days", tr)

	tr, err = localizer.lookupSingularTranslation("unknown", NoCtx, "%d day", 1)
	assert.Error(t, err)
	assert.Equal(t, "1 day", tr)

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

	assert.Equal(t, "ID", localizer.DGet("a", "id"))
	assert.Equal(t, "ID", localizer.DGetf("a", "id"))
	assert.Equal(t, "en_id", localizer.DGetf("a", "en_id"))
	assert.Equal(t, "id", localizer.DGetf("unknown", "id"))

	assert.Equal(t, "%d Tag", localizer.NGet("%d day", "%d days", 1))
	assert.Equal(t, "1 Tag", localizer.NGetf("%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.NGetf("%d car", "%d cars", 10, 10))

	assert.Equal(t, "1 Tag", localizer.DNGetf("a", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNGetf("a", "%d car", "%d cars", 10, 10))
	assert.Equal(t, "%d day", localizer.DNGet("unknown", "%d day", "%d days", 1))
	assert.Equal(t, "1 day", localizer.DNGetf("unknown", "%d day", "%d days", 1, 1))
	assert.Equal(t, "10 cars", localizer.DNGetf("unknown", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "Test mit Context", localizer.PGet("context", "Test with context"))
	assert.Equal(t, "Test mit Context", localizer.PGetf("context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.PGetf("other", "Test with context", 5))

	assert.Equal(t, "Test mit Context", localizer.DPGetf("a", "context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.DPGetf("a", "other", "Test with context", 5))
	assert.Equal(t, "Test with context", localizer.DPGetf("unknown", "context", "Test with context"))
	assert.Equal(t, "Test with context", localizer.DPGetf("unknown", "other", "Test with context", 5))

	assert.Equal(t, "%d Ergebniss", localizer.NPGetf("context", "%d result", "%d results", 1))
	assert.Equal(t, "1 Ergebniss", localizer.NPGetf("context", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 Ergebnisse", localizer.NPGetf("context", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 result", localizer.NPGetf("other", "%d result", "%d results", 1, 1))
	assert.Equal(t, "10 results", localizer.NPGetf("other", "%d result", "%d results", 10, 10))
	assert.Equal(t, "1 car", localizer.NPGetf("context", "%d car", "%d cars", 1, 1))
	assert.Equal(t, "10 cars", localizer.NPGetf("context", "%d car", "%d cars", 10, 10))

	assert.Equal(t, "%d Ergebniss", localizer.DNPGet("a", "context", "%d result", "%d results", 1))
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
			got, err := localizer.lookupSingularTranslation("a", tt.ctx, tt.msgID, tt.params...)
			tt.wantErr(t, err, fmt.Sprintf("pGettextErr(%q, %v)", tt.ctx, tt.msgID))
			assert.Equalf(t, tt.translated, got, "pGettextErr(%q, %v)", tt.ctx, tt.msgID)
		}
	})

	t.Run("translate plural", func(t *testing.T) {
		for idx, tt := range pluralTestData {
			got, err := localizer.lookupPluralTranslation("a", tt.ctx, tt.msgID, tt.plural, tt.n, tt.params...)
			if tt.wantErr(t, err, fmt.Sprintf("localizer idx=%d npGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)) {
				assert.Equalf(t, tt.translated, got, "localizer idx=%d npGettextErr(%q, %q)", idx, tt.ctx, tt.msgID)
			}
		}
	})

	t.Run("default text is returned for a non-existent domain", func(t *testing.T) {
		tr, err := localizer.lookupPluralTranslation("unknown", NoCtx, "%d day", "%d days", 1, 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)

		tr, err = localizer.lookupPluralTranslation("unknown", NoCtx, "%d day", "%d days", 10, 10)
		assert.Error(t, err)
		assert.Equal(t, "10 days", tr)

		tr, err = localizer.lookupSingularTranslation("unknown", NoCtx, "%d day", 1)
		assert.Error(t, err)
		assert.Equal(t, "1 day", tr)
	})
}

func TestLocalizer_LocalizeError(t *testing.T) {
	localizer := getLocalizerForTest(t)

	t.Run("error translation", func(t *testing.T) {
		want := "fehler"
		get := localizer.LocalizeError(errors.New("failure"))
		assert.Error(t, get)
		assert.IsType(t, &localize.Error{}, get)
		assert.Equal(t, want, get.Error())
	})

	t.Run("localizable error", func(t *testing.T) {
		raw := &testLocalizeErr{
			singular:  "failure",
			errorText: "other",
		}

		get := localizer.LocalizeError(raw)
		assert.Error(t, get)
		assert.IsType(t, raw, get)
		assert.Equal(t, raw.errorText, get.Error())

		raw.context = "errors"
		get = localizer.LocalizeError(raw)
		assert.Error(t, get)
		assert.IsType(t, &localize.Error{}, get)
		want := "fehler"
		assert.Equal(t, want, get.Error())

		raw.singular = "The kindergarten"
		raw.context = ""
		raw.domain = "z"
		raw.hasDomain = true
		get = localizer.LocalizeError(raw)
		assert.True(t, localizer.HasLocale())
		assert.Contains(t, localizer.Domains(), "z")
		assert.Error(t, get)
		if assert.IsType(t, &localize.Error{}, get) {
			loErr := get.(*localize.Error)
			want = "Der Kindergarten"
			assert.Equal(t, want, loErr.Translation)
		}
	})

	t.Run("no translation", func(t *testing.T) {
		err := errors.New("unknown error")
		get := localizer.LocalizeError(err)
		assert.Exactly(t, err, get)
	})
}

func TestKeyValueLocalizer(t *testing.T) {
	altDomain := "altjson"
	bundle, errB := NewBundle(
		WithFallbackLanguage(language.English),
		WithDefaultDomain("json"),
		WithDomainPath("json", testTranslationDir),
		WithDomainPath(altDomain, testTranslationDir),
		WithRequiredLanguage(language.German),
	)
	require.NoError(t, errB)
	require.NotNil(t, bundle)

	deLoc := NewKeyLocalizer(bundle, "de")
	enLoc := NewKeyLocalizer(bundle, "en")
	assert.Equal(t, language.German, deLoc.Language())

	t.Run("Uses fallback for unknown language", func(t *testing.T) {
		localizer := NewKeyLocalizer(bundle, "es")
		require.False(t, localizer.HasLocale())

		assert.Equal(t, "TODO List", localizer.Get("app.name"))
		assert.Equal(t, "I have a dog", localizer.NPGet("my-animals", "animal.dog", 1))
	})

	t.Run("Uses fallback when no translations are available", func(t *testing.T) {
		assert.Equal(t, "Monday", deLoc.Get("date.monday"))
		assert.Equal(t, "I do not have a cat", deLoc.NGet("animal.cat", 1))
		assert.Equal(t, "I do not have cats", deLoc.Get("animal.cat"))
	})

	t.Run("Fallback language is found", func(t *testing.T) {
		assert.True(t, enLoc.HasLocale())
	})

	t.Run("Key is output for unknown keys", func(t *testing.T) {
		assert.Equal(t, "unknown.key", enLoc.Get("unknown.key"))
	})

	t.Run("test translation functions", func(t *testing.T) {
		assert.Equal(t, "Tschau Max", deLoc.Getf("user.goodbye", "Max"))
		assert.Equal(t, "Willkommen Max", deLoc.DGetf(altDomain, "user.welcome", "Max"))
		assert.Equal(t, "Alternativer Domain", deLoc.DGet(altDomain, "domain.name"))
		assert.Equal(t, "I have 1 animal", enLoc.NGetf("animals", 1, 1))
		assert.Equal(t, "1 Tag", deLoc.DNGetf(altDomain, "days", 1, 1))
		assert.Equal(t, "weeks", enLoc.DNGet(altDomain, "week", 2))

		assert.Equal(t, "I have dogs", enLoc.PGet("my-animals", "animal.dog"))
		assert.Equal(t, "My pet's name is Lino", enLoc.PGetf("my-animals", "animal.name", "Lino"))
		assert.Equal(t, "I have a dog", enLoc.NPGet("my-animals", "animal.dog", 1))
		assert.Equal(t, "There are 100 animals in the world", enLoc.NPGetf("world", "animal", 100, 100))

		assert.Equal(t, "Test mit Context", deLoc.DPGet(altDomain, "context", "test"))
		assert.Equal(t, "Howdy partner", enLoc.DPGetf(altDomain, "cowboy", "user.welcome", "partner"))
		assert.Equal(t, "%d Ergebnisse", deLoc.DNPGet(altDomain, "context", "result", 2))
		assert.Equal(t, "2 Ergebnisse", deLoc.DNPGetf(altDomain, "context", "result", 2, 2))
	})
}
