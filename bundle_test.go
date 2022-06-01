package spreak

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog"
	"github.com/vorlif/spreak/internal/util"
)

func TestNewBundle(t *testing.T) {
	t.Run("error is returned when a nil option is passed", func(t *testing.T) {
		bundle, err := NewBundle(WithDomainPath(NoDomain, testdataStructureDir), nil)
		require.Error(t, err)
		require.Nil(t, bundle)
	})
}

func TestBundle_Domains(t *testing.T) {
	bundle, err := NewBundle(
		WithDomainPath(NoDomain, testdataStructureDir),
		WithDomainPath("a", testdataStructureDir),
		WithDomainPath("b", testdataStructureDir),
	)
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, 0, len(bundle.Domains()))
	assert.Equal(t, 0, len(bundle.SupportedLanguages()))
}

func TestBundle_SupportedLanguages(t *testing.T) {
	bundle, err := NewBundle(
		WithDomainPath(NoDomain, testdataStructureDir),
		WithRequiredLanguage(language.English, language.MustParse("de-at")),
		WithLanguage(language.Afrikaans),
	)
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, 1, len(bundle.Domains()))
	assert.Equal(t, 2, len(bundle.SupportedLanguages()))
}

func TestWithDefaultDomain(t *testing.T) {
	domainName := "my_domain"
	bundle, err := NewBundle(WithDefaultDomain(domainName))
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, domainName, bundle.defaultDomain)
}

func TestWithDomainFs(t *testing.T) {
	bundle, err := NewBundle(
		WithFallbackLanguage("en"),
		WithDomainPath("b", testdataStructureDir),
	)
	require.NoError(t, err)
	require.NotNil(t, bundle)

	assert.Equal(t, "", bundle.defaultDomain)
	assert.Equal(t, 1, len(bundle.locales))
}

func TestWithFallbackLanguage(t *testing.T) {
	bundle, err := NewBundle(
		WithFallbackLanguage("en"),
		WithDomainPath("b", testdataStructureDir),
	)
	if assert.NoError(t, err) {
		assert.NotNil(t, bundle)
		assert.Equal(t, language.English, bundle.fallbackLocale.language)
	}

	missingCount := 0
	missingCallback := func(err error) {
		missingCount++
	}

	bundle, err = NewBundle(
		WithFallbackLanguage("de_AT"),
		WithDomainPath(NoDomain, testdataStructureDir),
		WithMissingTranslationCallback(missingCallback),
	)

	if assert.NoError(t, err) {
		assert.NotNil(t, bundle)
		assert.Equal(t, language.MustParse("de_AT"), bundle.fallbackLocale.language)

		require.Equal(t, 0, missingCount)
		if !assert.Equal(t, 1, len(bundle.languages)) {
			for _, lang := range bundle.languages {
				t.Log(lang)
			}
		}
	}

	t.Run("error on invalid language", func(t *testing.T) {
		bundle, err = NewBundle(WithFallbackLanguage(1))
		assert.Error(t, err)
		assert.Nil(t, bundle)
	})

}

func TestWithLanguage(t *testing.T) {
	t.Run("Existing language is loaded", func(t *testing.T) {
		bundle, err := NewBundle(WithRequiredLanguage(language.German), WithDomainPath("a", testdataStructureDir))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		assert.True(t, bundle.CanLocalize())
		assert.Equal(t, 1, len(bundle.languages))
		assert.True(t, bundle.IsLanguageSupported(language.German))
		assert.False(t, bundle.IsLanguageSupported(language.French))
	})

	t.Run("Non-existent language leads to an error", func(t *testing.T) {
		bundle, err := NewBundle(WithRequiredLanguage(language.French), WithDomainPath("a", testdataStructureDir))
		require.Error(t, err)
		require.Nil(t, bundle)
	})

	t.Run("An optional language does not lead to any error", func(t *testing.T) {
		bundle, err := NewBundle(WithLanguage(language.French), WithDomainPath("a", testdataStructureDir))
		require.NoError(t, err)
		require.NotNil(t, bundle)

		assert.False(t, bundle.CanLocalize())
		assert.Equal(t, 0, len(bundle.languages))
		assert.False(t, bundle.IsLanguageSupported(language.German))
		assert.False(t, bundle.IsLanguageSupported(language.French))
	})
}

func TestWithLoader(t *testing.T) {
	t.Run("Nil is not a valid loader", func(t *testing.T) {
		bundle, err := NewBundle(WithDomainLoader(NoDomain, nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})

	t.Run("A domain can have only one loader", func(t *testing.T) {
		bundle, err := NewBundle(WithDomainPath(NoDomain, "/tmp"), WithDomainPath(NoDomain, testdataStructureDir))
		require.Error(t, err)
		require.Nil(t, bundle)

		loader := testLoader{
			f: func(lang language.Tag, domain string) (catalog.Catalog, error) {
				return nil, errors.New("test loader")
			},
		}
		bundle, err = NewBundle(WithDomainLoader(NoDomain, &loader), WithDomainLoader(NoDomain, &loader))
		require.Error(t, err)
		require.Nil(t, bundle)
	})

	t.Run("The passed loader is set", func(t *testing.T) {
		loader := &testLoader{f: func(lang language.Tag, domain string) (catalog.Catalog, error) {
			return nil, errors.New("not used")
		}}
		bundle, err := NewBundle(WithDomainLoader(NoDomain, loader))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		// TODO
	})
}

func TestWithMissingTranslationCallback(t *testing.T) {
	t.Run("Nil is a valid missing callback", func(t *testing.T) {
		bundle, err := NewBundle(WithMissingTranslationCallback(nil))
		require.NoError(t, err)
		require.NotNil(t, bundle)
	})

	t.Run("The passed missing callback is set", func(t *testing.T) {
		executionCount := 0
		var callback MissingTranslationCallback = func(err error) {
			executionCount++
		}

		bundle, err := NewBundle(WithMissingTranslationCallback(callback))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		assert.Equal(t, 0, executionCount)
		bundle.missingCallback(errors.New(NoDomain))
		assert.Equal(t, 1, executionCount)
	})
}

func TestWithPrintFuncGenerator(t *testing.T) {
	t.Run("The passed print function generator is set", func(t *testing.T) {
		executionCount := 0
		printer := &testPrinter{func(lang language.Tag) PrintFunc {
			executionCount++
			return func(str string, vars ...interface{}) string {
				return lang.String()
			}
		}}
		bundle, err := NewBundle(WithPrinter(printer))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		require.NotNil(t, bundle.printer)
		assert.Equal(t, 1, executionCount)

		printF := bundle.printer.GetPrintFunc(language.Und)
		result := printF("test")
		assert.Equal(t, language.Und.String(), result)
		assert.Equal(t, 2, executionCount)
	})

	t.Run("Nil is not a valid print function generator", func(t *testing.T) {
		bundle, err := NewBundle(WithPrinter(nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})
}

func TestWithPrintFunction(t *testing.T) {
	t.Run("The passed print functions is set", func(t *testing.T) {
		executionCount := 0
		var printF PrintFunc = func(str string, vars ...interface{}) string {
			executionCount++
			return str
		}
		bundle, err := NewBundle(WithPrintFunction(printF))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		require.NotNil(t, bundle.printer)
		assert.Equal(t, 0, executionCount)
		executionResult := bundle.printer.GetPrintFunc(language.Und)("test")
		assert.Equal(t, "test", executionResult)
		assert.Equal(t, 1, executionCount)
	})

	t.Run("Nil is not a valid print function", func(t *testing.T) {
		bundle, err := NewBundle(WithPrintFunction(nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})
}

func TestWithDomainFs1(t *testing.T) {
	t.Run("Nil is not a valid filesystem", func(t *testing.T) {
		bundle, err := NewBundle(WithDomainFs(NoDomain, nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})

	t.Run("The passed filesystem is set", func(t *testing.T) {
		fsys := util.DirFS(testdataStructureDir)
		bundle, err := NewBundle(WithDomainFs(NoDomain, fsys))
		require.NoError(t, err)
		require.NotNil(t, bundle)

		// TODO
	})
}

func TestWithFilesystemLoader(t *testing.T) {
	t.Run("Nil is not an valid option", func(t *testing.T) {
		bundle, err := NewBundle(WithFilesystemLoader(NoDomain, nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})
}

func TestWithErrorContext(t *testing.T) {
	t.Run("Empty string is set", func(t *testing.T) {
		bundle, err := NewBundle(WithErrorContext(NoCtx))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		assert.Equal(t, NoCtx, bundle.errContext)
	})

	t.Run("Value is set", func(t *testing.T) {
		bundle, err := NewBundle(WithErrorContext("my-context"))
		require.NoError(t, err)
		require.NotNil(t, bundle)
		assert.Equal(t, "my-context", bundle.errContext)
	})
}

func TestWithLanguageMatcherBuilder(t *testing.T) {
	t.Run("Nil is not an valid option", func(t *testing.T) {
		bundle, err := NewBundle(WithLanguageMatcherBuilder(nil))
		require.Error(t, err)
		require.Nil(t, bundle)
	})

	t.Run("language match-builder is set", func(t *testing.T) {
		var matcher language.Matcher
		builder := func(t []language.Tag, options ...language.MatchOption) language.Matcher {
			matcher = language.NewMatcher(t, options...)
			return matcher
		}
		bundle, err := NewBundle(WithLanguageMatcherBuilder(builder))
		require.NoError(t, err)
		require.NotNil(t, bundle)

		assert.Equal(t, matcher, bundle.languageMatcher)
	})
}
