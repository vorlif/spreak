package humanize

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/internal/util"
)

var testdataDir = filepath.FromSlash("../testdata/humanize")
var deLocaleDir = filepath.FromSlash("./locale/de")
var esLocaleDir = filepath.FromSlash("./locale/es")

var testGermanLocaleData = &LocaleData{
	Lang: language.German,
	Fs:   util.DirFS(deLocaleDir),
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

func createNewParcel(_ *testing.T) *Parcel {
	es := &LocaleData{Lang: language.Spanish, Fs: util.DirFS(esLocaleDir)}

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
		parcel, err := New(WithBundleOption(spreak.WithDomainPath(djangoDomain, "../testdata/humanize")))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "django")
		assert.Nil(t, parcel)
	})

	t.Run("test MuseNew with error panics", func(t *testing.T) {
		f := func() {
			_ = MustNew(WithBundleOption(spreak.WithDomainPath(djangoDomain, "../testdata/humanize")))
		}
		assert.Panics(t, f)
	})
}
