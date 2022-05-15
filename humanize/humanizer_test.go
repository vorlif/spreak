package humanize

import (
	"path/filepath"
	"testing"

	"golang.org/x/text/language"

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

func createNewParcel(t *testing.T) *Parcel {
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
