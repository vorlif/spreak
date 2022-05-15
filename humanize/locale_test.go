package humanize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLoader(t *testing.T) {
	t.Run("no languages registered", func(t *testing.T) {
		l := newLoader([]*LocaleData{})
		assert.Empty(t, l.locales)

		catl, err := l.Load(language.German, djangoDomain)
		assert.Error(t, err)
		assert.Nil(t, catl)
	})

	t.Run("languages registered", func(t *testing.T) {
		l := newLoader([]*LocaleData{testGermanLocaleData})
		assert.Len(t, l.locales, 1)

		catl, err := l.Load(language.German, djangoDomain)
		assert.NoError(t, err)
		assert.NotNil(t, catl)
	})

	t.Run("invalid filesystem", func(t *testing.T) {
		path := filepath.Join(testdataDir, "en")

		l := newLoader([]*LocaleData{{Lang: language.German, Fs: os.DirFS(path), Format: nil}})
		catl, err := l.Load(language.German, djangoDomain)
		assert.Error(t, err)
		assert.Nil(t, catl)
	})

	t.Run("invalid file", func(t *testing.T) {
		path := filepath.Join(testdataDir, "es")
		l := newLoader([]*LocaleData{{Lang: language.Spanish, Fs: os.DirFS(path), Format: nil}})
		catl, err := l.Load(language.Spanish, djangoDomain)
		assert.Error(t, err)
		assert.Nil(t, catl)
	})
}
