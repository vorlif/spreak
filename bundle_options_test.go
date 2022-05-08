package spreak

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/util"
)

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
	bundle, err := NewBundle(WithFallbackLanguage("en"))
	if assert.NoError(t, err) {
		assert.NotNil(t, bundle)
		assert.Equal(t, language.English, bundle.fallbackLanguage)
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
		assert.Equal(t, language.MustParse("de_AT"), bundle.fallbackLanguage)

		require.Equal(t, 0, missingCount)
		if !assert.Equal(t, 1, len(bundle.languages)) {
			for _, lang := range bundle.languages {
				t.Log(lang)
			}
		}
	}
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
	})

	t.Run("The passed loader is set", func(t *testing.T) {
		loader := &testLoader{f: func(lang language.Tag, domain string) (Catalog, error) {
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
		assert.Equal(t, 0, executionCount)

		printF := bundle.printer.GetPrintFunc(language.Und)
		result := printF("test")
		assert.Equal(t, language.Und.String(), result)
		assert.Equal(t, 1, executionCount)
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
