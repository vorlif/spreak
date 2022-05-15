package enGB

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("en-GB"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j M Y",
			TimeFormat:          "P",
			DateTimeFormat:      "j M Y, P",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "d/m/Y",
			ShortDatetimeFormat: "d/m/Y P",
			FirstDayOfWeek:      1,
		},
	}
}
