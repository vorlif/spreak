package th

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("th"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j F Y",
			TimeFormat:          "G:i",
			DateTimeFormat:      "j F Y, G:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "j M Y",
			ShortDatetimeFormat: "j M Y, G:i",
			FirstDayOfWeek:      0,
		},
	}
}
