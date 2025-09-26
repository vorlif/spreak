package humanize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

var testdataDir = filepath.FromSlash("../testdata/humanize")
var deLocaleDir = filepath.FromSlash("./locale/de")
var esLocaleDir = filepath.FromSlash("./locale/es")

var testGermanLocaleData = &LocaleData{
	Lang: language.German,
	Fs:   os.DirFS(deLocaleDir),
	Format: &FormatData{
		DateFormat:          "j. F Y",
		TimeFormat:          "H:i",
		DateTimeFormat:      "j. F Y H:i",
		YearMonthFormat:     "F Y",
		MonthDayFormat:      "j. F",
		ShortDateFormat:     "d.m.Y",
		ShortDatetimeFormat: "d.m.Y H:i",
		FirstDayOfWeek:      1,
	},
}

func createNewParcel(_ *testing.T) *Collection {
	es := &LocaleData{Lang: language.Spanish, Fs: os.DirFS(esLocaleDir)}

	return MustNew(WithLocale(testGermanLocaleData, es))
}

func createGermanHumanizer(t *testing.T) *Humanizer {
	p := createNewParcel(t)
	return p.CreateHumanizer(language.German)
}

func createSourceHumanizer(t *testing.T) *Humanizer {
	p := createNewParcel(t)
	return p.CreateHumanizer(language.English)
}

func TestNew(t *testing.T) {
	t.Run("test same domain returns error", func(t *testing.T) {
		collection, err := New(WithBundleOption(spreak.WithDomainPath(djangoDomain, "../testdata/humanize")))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "django")
		assert.Nil(t, collection)
	})

	t.Run("test MuseNew with error panics", func(t *testing.T) {
		f := func() {
			_ = MustNew(WithBundleOption(spreak.WithDomainPath(djangoDomain, "../testdata/humanize")))
		}
		assert.Panics(t, f)
	})
}
